package main

import (
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Object struct {
	Title, Body, Folder, FileTitle string
	Data                           []interface{}
	Err                            error
	Info                           string
}

type Objects struct {
	Title string
	Data  []string
	Err   error
	Info  string
}

type HomePage struct {
	Title, Body, Folder, FileTitle string
	Data                           []interface{}
	Atts                           []os.FileInfo
	Objs                           *Objects
	Err                            error
	Info                           string
}

func loadHomePage() (*Object, error) {
	o, err := NewObject("Home")
	if err != nil {
		return nil, err
	}
	o, err = load(o)
	if err != nil {
		return nil, err
	}
	o.Body = innerLink(o.Body)
	// list home attachments
	files, err := walk(o)
	if err != nil {
		return nil, err
	}
	Atts := []os.FileInfo{}
	for _, file := range files {
		Atts = append(Atts, file)
	}
	// list objects
	Objs, err := listObjects()
	if err != nil {
		return nil, err
	}
	o.Data = append(o.Data, Atts, Objs)
	return o, nil
}

func NewObject(title string) (*Object, error) {
	title, err := url.QueryUnescape(title)
	if err != nil {
		return nil, err
	}
	p := &Object{Title: title}
	p.Folder = filepath.Join(configs.dataPath, title)
	p.FileTitle = filepath.Join(p.Folder, title+".md")
	return p, nil
}

// save write done Body after NewObject() generate the p
func save(o *Object) error {
	if _, err := os.Stat(o.Folder); err != nil && os.IsNotExist(err) {
		os.MkdirAll(o.Folder, 0755)
	}
	return ioutil.WriteFile(o.FileTitle, []byte(o.Body), 0600)
}

// load read person info after NewObject() generate the p
func load(o *Object) (*Object, error) {
	body, err := ioutil.ReadFile(o.FileTitle)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return o, err
		}
		return nil, err
	}
	o.Body = string(body)
	return o, nil
}

func listObjects() (*Objects, error) {
	objs := &Objects{Title: "Objects list"}
	dirs, err := os.ReadDir(configs.dataPath)
	if err != nil {
		return nil, fmt.Errorf("error walking the path %q: %v\n", configs.dataPath, err)
	}
	for _, dir := range dirs {
		if dir.IsDir() && strings.ToLower(dir.Name()) != "home" {
			objs.Data = append(objs.Data, dir.Name())
		}
	}
	return objs, nil
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
	files := []fs.FileInfo{}
	err := filepath.Walk(o.Folder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if !info.IsDir() && filepath.Ext(path) != ".md" && info.Name()[:1] != "." {
			files = append(files, info)
		}

		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", o.Folder, err)
		return nil, err
	}
	return files, nil
}
