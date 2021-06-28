package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hi20160616/gears"
	"github.com/hi20160616/obox/configs"
	"github.com/hi20160616/obox/internal/data"
)

var dbFolderPath = filepath.Join(configs.Data.RootPath, configs.Data.DataPath)
var dbFilePath = filepath.Join(configs.Data.RootPath, "data.db")
var dbBakupFilePath = filepath.Join(configs.Data.RootPath, "data.db.bak")

func openData() error {
	return data.UnzipFiles(dbFilePath, dbFolderPath, configs.Data.Password)
}

func openBakData() error {
	return data.UnzipFiles(dbBakupFilePath, dbFolderPath, configs.Data.Password)
}

func closeData() error {
	return data.ZipFiles(dbFolderPath, dbFilePath, configs.Data.Password)
}

func backup() error {
	f, err := os.ReadFile(dbFilePath)
	if err != nil {
		return err
	}
	return os.WriteFile(dbBakupFilePath, f, os.ModePerm)
}

func main() {
	var err error
	defer func() {
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	}()

	if gears.Exists(dbFolderPath) {
		if err = closeData(); err != nil {
			return
		}
		if err = backup(); err != nil {
			return
		}
		if err = os.RemoveAll(dbFolderPath); err != nil {
			return
		}
		return
	} else if gears.Exists(dbFilePath) {
		if err = openData(); err != nil {
			return
		}
		if err = os.Remove(dbFilePath); err != nil {
			return
		}
		return
	} else if gears.Exists(dbBakupFilePath) {
		if err = data.UnzipFiles(
			dbBakupFilePath, dbFolderPath, configs.Data.Password); err != nil {
			return
		}
	}
}
