package main

import (
	"fmt"
	"io"
	"os"
	"mime/multipart"
	"net/http"
	"html/template"
	"encoding/hex"
	"crypto/sha1"
	"path/filepath"
)

// is this good?
var tpl = template.Must(template.ParseFiles("index.html"))

func main() {
	mux := buildMux()

	http.ListenAndServe(":8000", mux)
}

func buildMux() *http.ServeMux{
	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/post", createPostHandler)

	return mux
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func createPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20) // 32MB upload limit
		file, header, err := r.FormFile("file")
		if err != nil {
			fmt.Println(err)
			return
		}
		// TODO: maybe propogate error up to get a sense if it actually saved
		name := saveFile(file, header)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(201)
		w.Write([]byte(name))
	} else if r.Method == "GET" {
		w.Write([]byte(fmt.Sprintf("get")))
	}
} 

func saveFile(file multipart.File, header *multipart.FileHeader) string {
	defer file.Close()
	var name = hashName(header.Filename) + filepath.Ext(header.Filename)
	f, err := os.OpenFile("./" + name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return "/error"
	}
	defer f.Close()
	io.Copy(f, file)
	return name
}

// requires opened file
func hashName(name string) string {
	h := sha1.New()

	h.Write([]byte(name))

	return hex.EncodeToString(h.Sum(nil))
}