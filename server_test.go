package main

import "testing"

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
