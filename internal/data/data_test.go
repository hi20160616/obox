package data

import (
	"fmt"
	"testing"
)

func TestLoadHomePage(t *testing.T) {
	o, err := LoadHomePage()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(o)
}
