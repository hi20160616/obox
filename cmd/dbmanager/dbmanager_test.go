package main

import (
	"log"
	"os"
	"testing"
)

func TestOpenData(t *testing.T) {
	if err := openData(); err != nil {
		t.Error(err)
	}
}

func TestCloseData(t *testing.T) {
	if err := closeData(); err != nil {
		t.Error(err)
	}
}

func TestMain(m *testing.M) {
	log.Println("Do stuff BEFORE the tests!")
	exitVal := m.Run()
	log.Println("Do stuff AFTER the tests!")

	os.Exit(exitVal)
}
