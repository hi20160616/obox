package configs

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var Data = &configuration{}

type configuration struct {
	RootPath, DataPath, TmplPath, Password, Address string
}

func init() {
	var root string
	var err error
	if strings.Contains(os.Args[0], ".test") {
		root = "../../" // for test dbmanager
	} else {
		root, err = os.Getwd()
		if err != nil {
			log.Printf("config Getwd: %#v", err)
		}
	}
	f, err := os.ReadFile(filepath.Join(root, "configs/configs.json"))
	if err != nil {
		log.Printf("config ReadFile: %#v", err)
	}
	if err = json.Unmarshal(f, Data); err != nil {
		log.Printf("config Unmarshal err: %#v", err)
	}
	Data.RootPath = root
}
