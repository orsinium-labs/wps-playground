package main

import (
	"fmt"

	"github.com/life4/gweb/web"
)

type Python struct {
	pyodide web.Value
	doc     web.Document
	output  web.HTMLElement
}

func (py Python) print(text string, cls string) {
	el := py.doc.CreateElement("div")
	el.Attribute("class").Set("alert alert-" + cls)
	el.SetText(text)
	py.output.Node().AppendChild(el.Node())
}

func (py Python) PrintIn(text string) {
	py.print(text, "secondary")
}

func (py Python) PrintOut(text string) {
	py.print(text, "success")
}

func (py Python) PrintErr(text string) {
	py.print(text, "danger")
}

func (py Python) Run(cmd string) string {
	return py.pyodide.Call("runPython", cmd).String()
}

func (py Python) RunAndPrint(cmd string) {
	py.PrintIn(cmd)
	result := py.Run(cmd)
	py.PrintOut(result)
}

func (py Python) Install(pkg string) bool {
	cmd := fmt.Sprintf("micropip.install('%s')", pkg)
	py.PrintIn(cmd)
	_, fail := py.pyodide.Call("runPython", cmd).Promise().Get()
	if fail.Truthy() {
		py.PrintErr(fail.String())
		return false
	}
	py.PrintOut(fmt.Sprint(pkg, " installed"))
	return true
}

func (py Python) Set(name string, text string) {
	py.pyodide.Get("globals").Set(name, text)
}

func (py Python) Clear() {
	py.output.SetText("")
}

func (py Python) InitMicroPip() bool {
	py.PrintIn("import micropip")
	_, fail := py.pyodide.Call("loadPackage", "micropip").Promise().Get()
	if fail.Truthy() {
		py.PrintErr(fail.String())
		return false
	}
	py.Run("import micropip")
	py.PrintOut("True")
	return true
}
