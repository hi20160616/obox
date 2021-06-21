package main

var configs = &configuration{
	dataPath: "data",
	tmplPath: "tmpl",
}

type configuration struct {
	dataPath, tmplPath string
}
