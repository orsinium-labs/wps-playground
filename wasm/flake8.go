package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"syscall/js"

	"github.com/life4/gweb/web"
)

type Flake8 struct {
	script string
	btn    web.HTMLElement
	conf   web.HTMLElement
	doc    web.Document
	win    web.Window
	editor web.Value
	py     *Python
}

type Violation struct {
	Code        string
	Description string
	Context     string
	Line        int
	Column      int
	Plugin      string
}

func NewFlake8(win web.Window, doc web.Document, editor web.Value, py *Python) Flake8 {
	scripts := NewScripts()
	script := scripts.ReadFlake8()
	return Flake8{
		script: script,
		btn:    doc.Element("py-lint"),
		conf:   doc.Element("py-config"),
		doc:    doc,
		win:    win,
		editor: editor,
		py:     py,
	}

}

func (fh *Flake8) Register() {
	fh.btn.SetInnerHTML("run flake8")
	fh.btn.Set("disabled", false)

	wrapped := func(this js.Value, args []js.Value) interface{} {
		fh.btn.SetInnerHTML("running...")
		fh.btn.Set("disabled", true)
		fh.Run()
		go fh.Register()
		return true
	}
	fh.btn.Call("addEventListener", "click", js.FuncOf(wrapped))
}

func (fh *Flake8) Run() {
	fh.py.Clear()
	fh.py.Set("text", fh.editor.Call("getValue").String())
	fh.py.Set("config", fh.conf.Text())
	fh.py.RunAndPrint(fh.script)

	fh.py.Clear()
	fh.py.RunAndPrint("code")

	cmd := "'\\n'.join(app.formatter._out)"
	fh.py.PrintIn(cmd)
	result := fh.py.Run(cmd)
	fh.py.PrintOut(result)

	if result == "" {
		return
	}

	// read violations
	violations := make([]Violation, 0)
	for _, line := range strings.Split(result, "\n") {
		v := Violation{}
		err := json.Unmarshal([]byte(line), &v)
		if err != nil {
			fh.py.PrintErr(err.Error())
			return
		}
		violations = append(violations, v)
	}

	fh.py.Clear()
	fh.table(violations)
}

func (fh *Flake8) table(violations []Violation) {
	table := fh.doc.CreateElement("table")
	table.Attribute("class").Set("table table-sm")

	thead := fh.doc.CreateElement("thead")
	table.Node().AppendChild(thead.Node())
	tr := fh.doc.CreateElement("tr")
	thead.Node().AppendChild(tr.Node())

	cols := []string{"code", "descr", "pos", "context"}
	for _, name := range cols {
		th := fh.doc.CreateElement("th")
		th.SetText(name)
		tr.Node().AppendChild(th.Node())
	}

	tbody := fh.doc.CreateElement("tbody")
	table.Node().AppendChild(tbody.Node())

	for _, vl := range violations {
		tr := fh.doc.CreateElement("tr")

		td := fh.doc.CreateElement("td")
		td.SetText(vl.Code)
		tr.Node().AppendChild(td.Node())

		td = fh.doc.CreateElement("td")
		td.SetText(vl.Description)
		tr.Node().AppendChild(td.Node())

		td = fh.doc.CreateElement("td")
		td.SetText(fmt.Sprintf("%d:%d", vl.Line, vl.Column))
		tr.Node().AppendChild(td.Node())

		td = fh.doc.CreateElement("td")
		code := fh.doc.CreateElement("code")
		code.Attribute("class").Set("python")
		code.SetText(vl.Context)
		td.Node().AppendChild(code.Node())
		tr.Node().AppendChild(td.Node())

		tbody.Node().AppendChild(tr.Node())
	}

	fh.py.output.Node().AppendChild(table.Node())

	fh.win.Call("highlight")
}
