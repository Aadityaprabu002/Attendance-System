package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func foo(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmp, _ := template.ParseFiles("dummy.html")
		tmp.Execute(w, nil)
		return
	}

	err := r.ParseMultipartForm(4096)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r.FormValue("image"))
	file, handler, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Uploaded file:%+v\n", handler.Filename)
	fmt.Printf("File size:%+v\n", handler.Size)
	fmt.Printf("MIME header:%+v\n", handler.Header)

	cur, _ := os.Getwd()
	tempFile, err := ioutil.TempFile(cur, "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)

	w.Write([]byte("Success!"))
}

func initRouter() {
	r := mux.NewRouter()
	r.HandleFunc("/", foo)
	log.Fatal(http.ListenAndServe(":4040", r))
}
func main() {
	initRouter()
}
