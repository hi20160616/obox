package server

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/hi20160616/obox/configs"
	"github.com/hi20160616/obox/internal/server/handler"
	"github.com/hi20160616/obox/tmpl"
	"golang.org/x/crypto/ssh/terminal"
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
	var (
		pwd string = ""
		err error
	)
	fmt.Print("[!] Enter password: ")
	if runtime.GOOS == "windows" {
		r := bufio.NewReader(os.Stdin)
		pwd, err = r.ReadString('\n')
		if err != nil {
			return err
		}
	} else {
		b, err := terminal.ReadPassword(0)
		if err != nil {
			return err
		}
		pwd = string(b)
	}
	if configs.Data.Password != strings.TrimSpace(pwd) {
		fmt.Println("\n[-] Invalid password!")
		return ValidPasswd()
	}
	fmt.Println("\n[+] Pass!")
	return nil
}

func GetHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.FS(tmpl.FS))))
	mux.Handle("/src/", http.StripPrefix("/src/", http.FileServer(http.Dir(configs.Data.DataPath))))
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
