package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/hi20160616/obox/internal/data"
	"github.com/hi20160616/obox/tmpl"
	"github.com/yuin/goldmark"
)

// var templates = template.Must(template.ParseFS(tmpl, filepath.Join(tmplPath, "*.html")))
var templates = template.New("")

func init() {
	templates.Funcs(template.FuncMap{
		"smartTime":   smartTime,
		"markdown":    markdown,
		"byteCountSI": byteCountSI,
		"plusOne":     plusOne,
	})
	templates = template.Must(templates.ParseFS(tmpl.FS, "*.html"))
}

func markdown(in string) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(in), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func Derive(w http.ResponseWriter, tmplName string, o *data.Object) {
	if err := templates.ExecuteTemplate(w, tmplName+".html", o); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%v", err)
	}
}

func DeriveList(w http.ResponseWriter, os *data.Objects) {
	if err := templates.ExecuteTemplate(w, "list.html", os); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%v", err)
	}
}

func DeriveHome(w http.ResponseWriter, hp *data.Object) {
	if err := templates.ExecuteTemplate(w, "home.html", hp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%v", err)
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
