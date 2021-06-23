package main

import (
	"io/fs"
	"log"
	"net/http"
)

func main() {
	fs1, err := fs.Sub(tmpl, configs.tmplPath+"/bootstrap")
	if err != nil {
		log.Println(err)
	}
	http.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.FS(fs1))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/upload/", makeHandler(uploadHandler))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
