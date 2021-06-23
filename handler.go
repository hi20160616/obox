package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

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

func saveHandler(w http.ResponseWriter, r *http.Request, p *Object) {
	p.Body = r.FormValue("body")

	// p := &Object{Title: title, Body: body}
	if err := save(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+p.Title, http.StatusFound)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	p, err := NewObject("FrontObject")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	viewHandler(w, r, p)
}

var validPath = regexp.MustCompile("^/(edit|save|view|upload)/(.+)$")

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

// walk2 is encapsulated walk, that append fileinfos to o.Data
func walk2(o *Object) (*Object, error) {
	files, err := walk(o)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		o.Data = append(o.Data, file)
	}
	return o, nil
}

// walk get all files info in o.Folder
func walk(o *Object) ([]fs.FileInfo, error) {
	subDirToSkip := "skip"
	files := []fs.FileInfo{}
	err := filepath.Walk(o.Folder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == subDirToSkip {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}
		files = append(files, info)
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", o.Folder, err)
		return nil, err
	}
	return files, nil
}
