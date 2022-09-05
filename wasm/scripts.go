package main

import (
	"embed"
	"io/ioutil"
	"log"
	"strings"
)

//go:embed include/*
var included embed.FS

type Scripts struct{}

func (sc *Scripts) Read(fname string) []byte {
	file, err := included.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return content
}

// Read default flake8 config
func (sc *Scripts) ReadConfig() string {
	return string(sc.Read("include/setup.cfg"))
}

// Read the script for running flake8
func (sc *Scripts) ReadFlake8() string {
	return string(sc.Read("include/flake8.py"))
}

// Read the example.py shown in the input box by default
func (sc *Scripts) ReadExample() string {
	return string(sc.Read("include/example.py"))
}

// Read the requirements.txt file (list of Python dependencies).
func (sc *Scripts) ReadDeps() []string {
	content := string(sc.Read("include/requirements.txt"))
	return strings.Split(content, "\n")
}

func NewScripts() Scripts {
	return Scripts{}
}
