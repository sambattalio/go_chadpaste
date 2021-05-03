package main

import (
	"fmt"
	"io"
	"os"
        "time"
	"math/rand"
        "mime/multipart"
	"net/http"
	"html/template"
	"path/filepath"
)

// this should be variable... once all 5 length fill up move to 6
var HASH_LENGTH = 5
// is this good?
var tpl = template.Must(template.ParseFiles("index.html"))

func main() {
	mux := buildMux()


	http.ListenAndServe(":8000", mux)
}

func buildMux() *http.ServeMux{
	mux := http.NewServeMux()

        fs := http.FileServer(http.Dir("static"))
        fs2 := http.FileServer(http.Dir("f"))
        mux.Handle("/static/", http.StripPrefix("/static/", fs))
        mux.Handle("/f/", http.StripPrefix("/f/", fs2))
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
	var name = genAndCheckNewURL() + filepath.Ext(header.Filename)
	// check collision
	if _, err := os.Stat("./f/" + name); err == nil {
		// path/to/whatever exists

	}
	f, err := os.OpenFile("./f/" + name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return "/error"
	}
	defer f.Close()
	io.Copy(f, file)
	return name
}


// calhoun.io inspo 
const urlAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234566789"

var seededRandomization *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomString() string {
    b := make([]byte, HASH_LENGTH)
    for i := range b {
        b[i] = urlAlphabet[seededRandomization.Intn(len(urlAlphabet))]
    }
    return string(b)
}

func genAndCheckNewURL() string {
        potString := randomString()

        // TODO check against database and try a few times....
        // potentially add logic to increase size ofstring if too many

        return potString
}
