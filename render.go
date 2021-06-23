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
	"time"

	"github.com/yuin/goldmark"
)

//go:embed tmpl
var tmpl embed.FS

// var templates = template.Must(template.ParseFS(tmpl, filepath.Join(tmplPath, "*.html")))
var templates = template.New("")

func init() {
	templates.Funcs(template.FuncMap{
		"smartTime":   smartTime,
		"markdown":    markdown,
		"byteCountSI": byteCountSI,
		"plusOne":     plusOne,
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

func smartTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func byteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func plusOne(x int) int {
	return x + 1
}
