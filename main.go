package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/hi20160616/gears"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}

func run() error {
	if err := validPasswd(); err != nil {
		return err
	}

	if err := openData(); err != nil {
		return err
	}

	addr := ":" + generatePort()
	fmt.Println("[*] Server running at", addr)

	prepareHandler()
	return http.ListenAndServe(addr, nil)
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

func generatePort() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(9999))
}

func prepareHandler() {
	fs1, err := fs.Sub(tmpl, configs.TmplPath+"/bootstrap")
	if err != nil {
		log.Println(err)
	}
	http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.FS(fs1))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/upload/", makeHandler(uploadHandler))
	http.HandleFunc("/del/", delHandler)
	http.HandleFunc("/list/", listHandler)
	http.HandleFunc("/new/", newHandler)
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
