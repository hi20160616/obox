package main

import (
	"errors"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
)

type Object struct {
	Name, Body, Folder, Filename string
	Data                         []interface{}
	Err                          error
}

func NewObject(name string) (*Object, error) {
	name, err := url.QueryUnescape(name)
	if err != nil {
		return nil, err
	}
	p := &Object{Name: name}
	p.Folder = filepath.Join(configs.dataPath, name)
	p.Filename = filepath.Join(p.Folder, name+".md")
	return p, nil
}

// save write done Body after NewObject() generate the p
func save(p *Object) error {
	if _, err := os.Stat(p.Folder); err != nil && os.IsNotExist(err) {
		os.MkdirAll(p.Folder, 0755)
	}
	return ioutil.WriteFile(p.Filename, []byte(p.Body), 0600)
}

// load read person info after NewObject() generate the p
func load(p *Object) (*Object, error) {
	body, err := ioutil.ReadFile(p.Filename)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return p, err
		}
		return nil, err
	}
	p.Body = string(body)
	return p, nil
}
