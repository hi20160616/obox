package main

import (
	"fmt"
	"testing"
)

func TestZipFiles2(t *testing.T) {
	in, out := configs.DataPath, "./data.zip"
	if err := zipFiles2(in, out, "password"); err != nil {
		t.Error(err)
	}

	unzipFiles("./data.zip", "./zipTest", "password")
}

func TestZipFiles(t *testing.T) {
	in, out := configs.DataPath, "./data.zip"
	fmt.Println(in)
	if err := zipFiles(in, out, "password"); err != nil {
		t.Error(err)
	}

	unzipFiles("./data.zip", "./zipTest", "password")
}
