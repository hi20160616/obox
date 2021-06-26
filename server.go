package main

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hi20160616/gears"
)

var addr string

type Server struct {
	http.Server
}

func NewServer(address string) (*Server, error) {

	return &Server{http.Server{
		Addr:    address,
		Handler: GetHandler(),
	}}, nil
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return ctx.Err()
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.Shutdown(context.Background()); err != nil {
		return err
	}
	return ctx.Err()
}

func validPasswd() error {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("[!] Enter password: ")
	pwd, err := r.ReadString('\n')
	if err != nil {
		return err
	}
	if configs.Password != strings.TrimSpace(pwd) {
		fmt.Println("[-] Invalid password!")
		return validPasswd()
	}
	fmt.Println("[+] Pass!")
	return nil
}

func GetHandler() *http.ServeMux {
	mux := http.NewServeMux()
	fs1, err := fs.Sub(tmpl, configs.TmplPath+"/bootstrap")
	if err != nil {
		log.Println(err)
	}
	mux.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.FS(fs1))))
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/view/", makeHandler(viewHandler))
	mux.HandleFunc("/edit/", makeHandler(editHandler))
	mux.HandleFunc("/save/", makeHandler(saveHandler))
	mux.HandleFunc("/upload/", makeHandler(uploadHandler))
	mux.HandleFunc("/del/", delHandler)
	mux.HandleFunc("/list/", listHandler)
	mux.HandleFunc("/new/", newHandler)
	return mux
}

func openData() error {
	if gears.Exists(filepath.Join(configs.RootPath, "data.db")) {
		return unzipFiles("data.db", configs.RootPath, configs.Password)
	}
	return nil
}

func closeData() error {
	if err := zipFiles(configs.DataPath,
		filepath.Join(configs.RootPath, "data.db"),
		configs.Password); err != nil {
		return err
	}
	return os.RemoveAll(filepath.Join(configs.RootPath, configs.DataPath))
}
