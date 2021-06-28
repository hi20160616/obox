package tmpl

import (
	"fmt"
	"testing"
)

func TestTmpl(t *testing.T) {
	fss, err := FS.Open("bootstrap/bootstrap.min.css")
	if err != nil {
		t.Error(err)
	}
	fi, err := fss.Stat()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fi)
}
