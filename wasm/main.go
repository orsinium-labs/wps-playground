package main

import (
	"encoding/base64"
	"fmt"

	"github.com/life4/gweb/web"
)

func main() {
	window := web.GetWindow()
	doc := window.Document()
	doc.SetTitle("FlakeHell online")

	// init code editor
	input := doc.Element("py-code")
	scripts := NewScripts()
	ex := scripts.ReadExample()
	input.SetInnerHTML(ex)
	editor := window.Get("CodeMirror").Call("fromTextArea",
		input,
		map[string]interface{}{
			"lineNumbers": true,
		})

	config := doc.Element("py-config")
	config.SetText(scripts.ReadConfig())

	// load python
	py := Python{doc: doc, output: doc.Element("py-output")}
	py.PrintIn("Load Python")
	window.Get("languagePluginLoader").Promise().Get()
	py.PrintOut("Python is ready")
	py.pyodide = window.Get("pyodide")
	py.RunAndPrint("'Hello world!'")

	ok := py.InitMicroPip()
	if !ok {
		return
	}

	// skip nighty packages
	skip := map[string]string{
		"flake8-quotes":         "2.1.2",
		"flake8-bugbear":        "19.3",
		"flake8-rst-docstrings": "0.0.12",
		"flake8-eradicate":      "0.3.0",
		"flake8-isort":          "3.0.1",
		"flake8-bandit":         "2.1.1",
		"lazy-object-proxy":     "1.4.0",
		"wrapt":                 "1.11.0",
		"typed-ast":             "1.4.1",
		"six":                   "1.15.0",
	}
	for pname, pversion := range skip {
		cmd := "micropip.PACKAGE_MANAGER.installed_packages['%s'] = '%s'"
		py.Run(fmt.Sprintf(cmd, pname, pversion))
	}

	// install dependencies
	py.Clear()
	py.Install("flake8==3.8.0")
	py.Install("setuptools")
	py.Install("entrypoints")
	py.Install("flake8-builtins==1.5.3")
	py.Install("flakehell==0.7.0")

	// install non-wheel dependencies
	py.RunAndPrint("import sys")
	py.RunAndPrint("sys.path.insert(0, '.')")
	unzip := []string{
		"/flake8_quotes.zip",
		"/lazy_object_proxy.zip",
		"/wrapt.zip",
		"/typed_ast.zip",
	}
	extract := scripts.ReadExtract()
	for _, name := range unzip {
		archive := scripts.Read(name)
		encoded := base64.StdEncoding.EncodeToString(archive)
		py.Set("archive", encoded)
		py.RunAndPrint(extract)
	}

	installer := Installer{py: &py, doc: doc, win: window}
	installer.Init()

	flakehell := NewFlakeHell(window, doc, editor, &py)
	flakehell.Register()

	py.Clear()
	py.PrintOut("Ready!")

	select {}
}
