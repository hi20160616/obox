package configs

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var Data = &configuration{}

type configuration struct {
	RootPath, DataPath, TmplPath, Password string
}

func init() {
	root, err := os.Getwd()
	if err != nil {
		log.Printf("config Getwd: %#v", err)
	}
	// root = "../../../" // for test handler
	f, err := os.ReadFile(filepath.Join(root, "configs/configs.json"))
	if err != nil {
		log.Printf("config ReadFile: %#v", err)
	}
	if err = json.Unmarshal(f, Data); err != nil {
		log.Printf("config Unmarshal err: %#v", err)
	}
	Data.RootPath = root
}
