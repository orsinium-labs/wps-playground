package main

import (
	"strings"
	"syscall/js"

	"github.com/life4/gweb/web"
)

var packages = [...]string{
	"deal==4.2.0",
	"pylint==2.6.0",
	"wemake-python-styleguide==0.14.1",
}
var links = [...]string{
	"https://github.com/life4/deal",
	"https://github.com/PyCQA/pylint/",
	"https://github.com/wemake-services/wemake-python-styleguide",
}

type Installer struct {
	py        *Python
	doc       web.Document
	win       web.Window
	installed map[string]bool
}

func (inst *Installer) Init() {
	inst.installed = make(map[string]bool, len(packages))
	table := inst.doc.CreateElement("table")
	table.Attribute("class").Set("table table-sm")

	target := inst.doc.Element("py-plugins")
	target.SetText("")
	target.Node().AppendChild(table.Node())

	tbody := inst.doc.CreateElement("tbody")
	table.Node().AppendChild(tbody.Node())

	for i, pkg := range packages {
		tr := inst.doc.CreateElement("tr")
		tbody.Node().AppendChild(tr.Node())

		a := inst.doc.CreateElement("a")
		a.Attribute("href").Set(links[i])
		a.Attribute("target").Set("_blank")
		a.SetText(strings.SplitN(pkg, "==", 2)[0])

		td := inst.doc.CreateElement("td")
		td.Node().AppendChild(a.Node())
		tr.Node().AppendChild(td.Node())

		btn := inst.doc.CreateElement("button")
		btn.Class().Set("btn btn-outline-primary")
		btn.SetText("install")
		inst.bound(pkg, btn)

		td = inst.doc.CreateElement("td")
		td.Node().AppendChild(btn.Node())
		tr.Node().AppendChild(td.Node())
	}
}

func (inst *Installer) bound(pkg string, btn web.HTMLElement) {
	handle := func() {
		btn.Set("disabled", true)
		btn.SetText("installing...")
		inst.py.Install(pkg)
		btn.Class().Set("btn btn-outline-light")
		btn.SetText("installed")
	}

	wrapped := func(this js.Value, args []js.Value) interface{} {
		go handle()
		return true
	}
	btn.Call("addEventListener", "click", js.FuncOf(wrapped))
}
