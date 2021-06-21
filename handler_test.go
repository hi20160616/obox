package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// https://golang.org/pkg/net/http/httptest/#example_Server
func TestEditHandler(t *testing.T) {
	handler := http.HandlerFunc(makeHandler(editHandler))

	req := httptest.NewRequest("GET", "http://localhost:8080/edit/TestObject", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
}
