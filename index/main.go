package main

import (
	home "attsys/home/backend"
	users "attsys/student/backend"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func initRouter() {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	r.HandleFunc("/", home.Homepage)
	r.HandleFunc("/signin", users.Signin)
	r.HandleFunc("/signup", users.Signup)
	r.HandleFunc("/matchface", users.MatchFace)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	initRouter()
}
