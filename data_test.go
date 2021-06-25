package main

import (
	"fmt"
	"testing"
)

func TestLoadHomePage(t *testing.T) {
	o, err := loadHomePage()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(o)
}
