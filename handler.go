package main

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var validPath = regexp.MustCompile("^/(edit|save|view|upload|del)/(.+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, *Object)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		o, err := NewObject(m[2])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fn(w, r, o)
	}
}

func newHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	o, err := NewObject(name)
	if err != nil {
		o.Err = err
		Derive(w, "new", o)
		return
	}
	Derive(w, "new", o)
}

func viewHandler(w http.ResponseWriter, r *http.Request, o *Object) {
	o, err := load(o)
	if err != nil {
		http.Redirect(w, r, "/edit/"+o.Title, http.StatusFound)
		return
	}
	o.Body = innerLink(o.Body)
	o, err = walk2(o)
	if err != nil {
		o.Err = err
		Derive(w, "view", o)
		return
	}
	Derive(w, "view", o)
}

func editHandler(w http.ResponseWriter, r *http.Request, o *Object) {
	if o.Err != nil {
		Derive(w, "edit", o)
		return
	}
	o, err := load(o)
	if err != nil {
		o.Err = err
		Derive(w, "edit", o)
		return
	}
	o, err = walk2(o)
	if err != nil {
		o.Err = err
		Derive(w, "edit", o)
		return
	}
	Derive(w, "edit", o)
}

func uploadHandler(w http.ResponseWriter, r *http.Request, o *Object) {
	upload := func() error {
		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		r.ParseMultipartForm(10 << 20)
		file, handler, err := r.FormFile("myFile")
		if err != nil {
			return errors.WithMessage(err, "Error Retrieving the File")
		}
		defer file.Close()
		f, err := os.Create(filepath.Join(o.Folder, handler.Filename))
		if err != nil {
			return errors.WithMessagef(err, "Cannot create file as %s",
				handler.Filename)
		}
		defer f.Close()

		// read all of the contents of uploaded files into a
		// byte array
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return errors.WithMessagef(err, "Read bytes from %s error", handler.Filename)
		}

		if _, err := f.Write(fileBytes); err != nil {
			return errors.WithMessagef(err, "Write bytes to file error")
		}

		o.Info = "File Successfully Uploaded!"
		return nil
	}

	if err := upload(); err != nil {
		o.Err = err
		Derive(w, "edit", o)
		return
	}

	o, err := walk2(o)
	if err != nil {
		o.Err = err
		Derive(w, "edit", o)
		return
	}
	Derive(w, "edit", o)
}

func delHandler(w http.ResponseWriter, r *http.Request) {
	ss := strings.Split(r.URL.Path, "/")
	if len(ss) >= 4 {
		if err := os.Remove(filepath.Join(configs.DataPath, ss[2], ss[3])); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	o, err := NewObject(ss[2])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	editHandler(w, r, o)
}

func saveHandler(w http.ResponseWriter, r *http.Request, o *Object) {
	o.Body = r.FormValue("body")

	// o := &Object{Title: title, Body: body}
	if err := save(o); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+o.Title, http.StatusFound)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	hp, err := loadHomePage()
	if err != nil {
		if os.IsNotExist(err) {
			http.Redirect(w, r, "edit/Home", http.StatusFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	DeriveHome(w, hp)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	objs, err := listObjects()
	if err != nil {
		objs.Err = err
		DeriveList(w, objs)
		return
	}
	DeriveList(w, objs)
}
