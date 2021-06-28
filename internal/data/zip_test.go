package data

import (
	"fmt"
	"github.com/hi20160616/obox/configs"
	"testing"
)

func TestZipFiles2(t *testing.T) {
	in, out := configs.Data.DataPath, "./data.zip"
	if err := zipFiles2(in, out, "password"); err != nil {
		t.Error(err)
	}

	UnzipFiles("./data.zip", "./zipTest", "password")
}

func TestZipFiles(t *testing.T) {
	in, out := configs.Data.DataPath, "./data.zip"
	fmt.Println(in)
	if err := ZipFiles(in, out, "password"); err != nil {
		t.Error(err)
	}

	UnzipFiles("./data.zip", "./zipTest", "password")
}
