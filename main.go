package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 50

func main() {
	http.HandleFunc("/upload", upload)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	err := r.ParseMultipartForm(MAX_UPLOAD_SIZE)

	if err != nil {
		log.Println(err)
		return
	}

	file, handler, err := r.FormFile("file")

	defer file.Close()
	if err != nil {
		log.Println(err)
		return
	}

	f, err := os.OpenFile("./up/"+handler.Filename+time.Now().String(), os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()

	if err != nil {
		log.Println(err)
		return
	}

	_, err = io.Copy(f, file)

	if err != nil {
		log.Println(err)
		return
	}
}
