package main

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/yuin/goldmark"
)

//go:embed tmpl
var tmpl embed.FS

// var templates = template.Must(template.ParseFS(tmpl, filepath.Join(tmplPath, "*.html")))
var templates = template.New("")

func init() {
	templates.Funcs(template.FuncMap{
		"markdown": markdown,
	})
	templates = template.Must(templates.ParseFS(tmpl, filepath.Join(configs.tmplPath, "*.html")))
}

// pattern like `[!foobar]` means a inter-page need to be made as link
var innerObject = regexp.MustCompile(`\[!.+\]`)

func innerLink(body string) string {
	repl := func(pagename string) string {
		pagename = pagename[2 : len(pagename)-1]
		origin := pagename
		pagename = strings.ReplaceAll(pagename, " ", "-")
		return fmt.Sprintf("[%s](/view/%s)", origin, pagename)
	}

	return innerObject.ReplaceAllStringFunc(body, repl)
}

func markdown(in string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(in), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func Derive(w http.ResponseWriter, tmpl string, p *Object) {
	if err := templates.ExecuteTemplate(w, tmpl+".html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("err template: %s.html\n\terror: %v", tmpl, err)
	}
}
