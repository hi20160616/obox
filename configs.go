package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var configs = &configuration{}

type configuration struct {
	RootPath, DataPath, TmplPath, Password string
}

func init() {
	root, err := os.Getwd()
	if err != nil {
		log.Printf("config Getwd: %#v", err)
	}
	// root = "../../../" // for test handler
	f, err := os.ReadFile(filepath.Join(root, "configs.json"))
	if err != nil {
		log.Printf("config ReadFile: %#v", err)
	}
	if err = json.Unmarshal(f, configs); err != nil {
		log.Printf("config Unmarshal err: %#v", err)
	}
	configs.RootPath = root
}
