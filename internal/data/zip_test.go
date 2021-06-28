package data

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/hi20160616/obox/configs"
)

var dbFolderPath = filepath.Join(configs.Data.RootPath, configs.Data.DataPath)
var dbFilePath = filepath.Join(configs.Data.RootPath, "data.zip")

func TestZipFiles(t *testing.T) {
	in, out := dbFolderPath, dbFilePath
	if err := ZipFiles(in, out, "password"); err != nil {
		t.Error(err)
	}

	// UnzipFiles(out, in, "password")
}

func TestUnZipFiles(t *testing.T) {
	out, in := dbFolderPath, dbFilePath
	if err := UnzipFiles(in, out, configs.Data.Password); err != nil {
		t.Error(err)
	}
}

func TestMain(m *testing.M) {
	log.Println("Do stuff BEFORE the tests!")
	exitVal := m.Run()
	log.Println("Do stuff AFTER the tests!")

	os.Exit(exitVal)
}
