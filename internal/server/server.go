package server

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/hi20160616/gears"
	"github.com/hi20160616/obox/configs"
	"github.com/hi20160616/obox/internal/data"
	"github.com/hi20160616/obox/internal/server/handler"
	"github.com/hi20160616/obox/tmpl"
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

func ValidPasswd() error {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("[!] Enter password: ")
	pwd, err := r.ReadString('\n')
	if err != nil {
		return err
	}
	if configs.Data.Password != strings.TrimSpace(pwd) {
		fmt.Println("[-] Invalid password!")
		return ValidPasswd()
	}
	fmt.Println("[+] Pass!")
	return nil
}

func GetHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.FS(tmpl.FS))))
	mux.HandleFunc("/", handler.HomeHandler)
	mux.HandleFunc("/view/", handler.MakeHandler(handler.ViewHandler))
	mux.HandleFunc("/edit/", handler.MakeHandler(handler.EditHandler))
	mux.HandleFunc("/save/", handler.MakeHandler(handler.SaveHandler))
	mux.HandleFunc("/upload/", handler.MakeHandler(handler.UploadHandler))
	mux.HandleFunc("/del/", handler.DelHandler)
	mux.HandleFunc("/list/", handler.ListHandler)
	mux.HandleFunc("/new/", handler.NewHandler)
	return mux
}

func openData() error {
	if gears.Exists(filepath.Join(configs.Data.RootPath, "data.db")) {
		return data.UnzipFiles("data.db", configs.Data.RootPath, configs.Data.Password)
	}
	return nil
}

func closeData() error {
	if err := data.ZipFiles(configs.Data.DataPath,
		filepath.Join(configs.Data.RootPath, "data.db"),
		configs.Data.Password); err != nil {
		return err
	}
	return os.RemoveAll(filepath.Join(configs.Data.RootPath, configs.Data.DataPath))
}