package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// https://golang.org/pkg/net/http/httptest/#example_Server
func TestEditHandler(t *testing.T) {
	handler := http.HandlerFunc(MakeHandler(EditHandler))

	req := httptest.NewRequest("GET", "http://localhost:8080/edit/TestObject", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
}

func TestDelHandler(t *testing.T) {
	handler := http.HandlerFunc(DelHandler)

	req := httptest.NewRequest("GET", "http://localhost:8080/del/FrontObject/1024px-WoW_icon.png", nil)
	w := httptest.NewRecorder()
	handler(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
}
