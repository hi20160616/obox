package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
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

	validPasswd()

	addr := ":" + generatePort()
	fmt.Println("[*] Server running at", addr)

	log.Fatal(http.ListenAndServe(addr, nil))
}

func validPasswd() error {
	r := bufio.NewReader(os.Stdin)
	fmt.Print("[!] Enter password: ")
	pwd, err := r.ReadString('\n')
	if err != nil {
		log.Println(err)
		os.Exit(1)
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
