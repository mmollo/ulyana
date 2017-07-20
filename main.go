package main

import (
	"flag"
	"github.com/nu7hatch/gouuid"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 50

var ud *string

func main() {

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	ud = flag.String("d", cwd+"/html/upload", "Upload directory")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir("html")))
	http.HandleFunc("/upload", handler_upload)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handler_upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" && r.Method != "PUT" {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	err := r.ParseMultipartForm(MAX_UPLOAD_SIZE)

	if err != nil {
		log.Println(err)
		return
	}

	file, handler, err := r.FormFile("file")
	mimetype := handler.Header.Get("Content-Type")

	switch mimetype {
	case "image/jpeg":
	case "image/png":
	case "image/gif":
		break
	default:
		return
	}

	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	uid := upload(file)
	if r.Method == "PUT" {
		w.Write([]byte(uid))
		return
	}

	if r.Method == "POST" {
		w.Header().Set("Location", "/upload/"+uid+".jpg")
		w.WriteHeader(http.StatusFound)
		return
	}
}

func upload(file multipart.File) string {
	u4, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
		return ""
	}

	filename := u4.String()
	f, err := os.OpenFile(*ud+"/"+filename+".jpg", os.O_WRONLY|os.O_CREATE, 0666)
	defer f.Close()

	if err != nil {
		log.Println(err)
		return ""
	}

	_, err = io.Copy(f, file)

	if err != nil {
		log.Println(err)
		return ""
	}

	return filename
}
