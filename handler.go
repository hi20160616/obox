package main

import (
	"net/http"
	"regexp"
)

func viewHandler(w http.ResponseWriter, r *http.Request, p *Object) {
	p, err := load(p)
	if err != nil {
		http.Redirect(w, r, "/edit/"+p.Name, http.StatusFound)
		return
	}
	p.Body = innerLink(p.Body)
	Derive(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, p *Object) {
	p, err := load(p)
	if err != nil {
		p.Err = err
		Derive(w, "edit", p)
	} else {
		Derive(w, "edit", p)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request, p *Object) {
	p.Body = r.FormValue("body")

	// p := &Object{Name: title, Body: body}
	if err := save(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+p.Name, http.StatusFound)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p, err := NewObject("FrontObject")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	viewHandler(w, r, p)
}

var validPath = regexp.MustCompile("^/(edit|save|view)/(.+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, *Object)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		p, err := NewObject(m[2])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fn(w, r, p)
	}
}
