package main

import (
	"github.com/life4/gweb/web"
)

func main() {
	window := web.GetWindow()
	doc := window.Document()
	doc.SetTitle("WPS Playground")

	// init code editor
	input := doc.Element("py-code")
	scripts := NewScripts()
	ex := scripts.ReadExample()
	input.SetInnerHTML(ex)
	doc.Element("py-config").SetInnerHTML(scripts.ReadConfig())
	editor := window.Get("CodeMirror").Call("fromTextArea",
		input.JSValue(),
		map[string]interface{}{
			"lineNumbers": true,
		},
	)

	// load python
	py := Python{doc: doc, output: doc.Element("py-output")}
	py.PrintIn("Loading Python...")
	var err web.Value
	py.pyodide, err = window.Call("loadPyodide").Promise().Get()
	if !err.IsUndefined() {
		py.PrintErr(err.String())
		return
	}
	py.PrintOut("Python is ready")
	py.RunAndPrint("'Hello world!'")

	ok := py.InitMicroPip()
	if !ok {
		return
	}

	// install dependencies
	py.Clear()
	for _, dep := range scripts.ReadDeps() {
		if dep == "" {
			continue
		}
		if dep[0] == '#' {
			continue
		}
		py.Install(dep)
	}

	flake8 := NewFlake8(window, doc, editor, &py)
	flake8.Register()

	py.Clear()
	py.PrintOut("Ready!")

	select {}
}
