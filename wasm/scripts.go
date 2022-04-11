package main

import (
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/life4/flakehell-online/wasm/statik"

	"github.com/rakyll/statik/fs"
)

type Scripts struct {
	sfs http.FileSystem
}

func (sc *Scripts) Read(fname string) []byte {
	file, err := sc.sfs.Open(fname)
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

func (sc *Scripts) ReadConfig() string {
	return string(sc.Read("/config.toml"))
}

func (sc *Scripts) ReadFlakeHell() string {
	return string(sc.Read("/flakehell.py"))
}

func (sc *Scripts) ReadExample() string {
	return string(sc.Read("/example.py"))
}

func (sc *Scripts) ReadExtract() string {
	return string(sc.Read("/extract.py"))
}

func NewScripts() Scripts {
	sfs, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	return Scripts{sfs: sfs}
}
