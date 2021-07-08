package handler

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hi20160616/gears"
	"github.com/hi20160616/obox/configs"
	"github.com/hi20160616/obox/internal/data"
	"github.com/hi20160616/obox/internal/server/render"
	"github.com/pkg/errors"
)

var validPath = regexp.MustCompile("^/(edit|save|view|upload|del)/(.+)$")

func MakeHandler(fn func(http.ResponseWriter, *http.Request, *data.Object)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		o, err := data.NewObject(m[2])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fn(w, r, o)
	}
}

func NewHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	o, err := data.NewObject(name)
	if err != nil {
		o.Err = err
		render.Derive(w, "new", o)
		return
	}
	render.Derive(w, "new", o)
}

func ViewHandler(w http.ResponseWriter, r *http.Request, o *data.Object) {
	o, err := data.Load(o)
	if err != nil {
		http.Redirect(w, r, "/edit/"+o.Title, http.StatusFound)
		return
	}
	o.Body = data.InnerLink(o.Body)
	o, err = data.Walk2(o)
	if err != nil {
		o.Err = err
		render.Derive(w, "view", o)
		return
	}
	render.Derive(w, "view", o)
}

func EditHandler(w http.ResponseWriter, r *http.Request, o *data.Object) {
	if o.Err != nil {
		render.Derive(w, "edit", o)
		return
	}
	o, err := data.Load(o)
	if err != nil {
		o.Err = err
		render.Derive(w, "edit", o)
		return
	}
	o, err = data.Walk2(o)
	if err != nil {
		o.Err = err
		render.Derive(w, "edit", o)
		return
	}
	render.Derive(w, "edit", o)
}

func UploadHandler(w http.ResponseWriter, r *http.Request, o *data.Object) {
	upload := func() error {
		// Parse our multipart form, 10 << 20 specifies a maximum
		// upload of 10 MB files.
		r.ParseMultipartForm(10 << 20)
		file, handler, err := r.FormFile("myFile")
		if err != nil {
			return errors.WithMessage(err, "Error Retrieving the File")
		}
		defer file.Close()
		saveFilePath := filepath.Join(o.Folder, handler.Filename)
		if !configs.Data.UploadOverwrite {
			if gears.Exists(saveFilePath) {
				saveFilePath = filepath.Join(o.Folder, "_"+handler.Filename)
			}
		}
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
		render.Derive(w, "edit", o)
		return
	}

	EditHandler(w, r, o)
}

func DelHandler(w http.ResponseWriter, r *http.Request) {
	ss := strings.Split(r.URL.Path, "/")
	if len(ss) >= 4 {
		if err := os.Remove(filepath.Join(configs.Data.DataPath, ss[2], ss[3])); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	o, err := data.NewObject(ss[2])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	EditHandler(w, r, o)
}

func SaveHandler(w http.ResponseWriter, r *http.Request, o *data.Object) {
	o.Body = r.FormValue("body")

	// o := &Object{Title: title, Body: body}
	if err := data.Save(o); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+o.Title, http.StatusFound)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	hp, err := data.LoadHomePage()
	if err != nil {
		if os.IsNotExist(err) {
			http.Redirect(w, r, "edit/Home", http.StatusFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	render.DeriveHome(w, hp)
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	objs, err := data.ListObjects()
	if err != nil {
		objs.Err = err
		render.DeriveList(w, objs)
		return
	}
	render.DeriveList(w, objs)
}
