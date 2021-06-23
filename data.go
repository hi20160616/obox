package main

import (
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
)

type Object struct {
	Title, Body, Folder, FileTitle string
	Data                           []interface{}
	Err                            error
	Info                           string
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
func save(p *Object) error {
	if _, err := os.Stat(p.Folder); err != nil && os.IsNotExist(err) {
		os.MkdirAll(p.Folder, 0755)
	}
	return ioutil.WriteFile(p.FileTitle, []byte(p.Body), 0600)
}

// load read person info after NewObject() generate the p
func load(p *Object) (*Object, error) {
	body, err := ioutil.ReadFile(p.FileTitle)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return p, err
		}
		return nil, err
	}
	p.Body = string(body)
	return p, nil
}
