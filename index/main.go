package main

import (
	home "attsys/home/backend"
	users "attsys/user/backend"

	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var DB gorm.DB

func initRouter() {
	r := mux.NewRouter()
	r.HandleFunc("/", home.Homepage)
	r.HandleFunc("/signin", users.Signin)
	r.HandleFunc("/signup", users.Signup)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	initRouter()
	// dsn := "host=localhost user=aaditya password=1234 dbname=postgres port=5432 sslmode=disable TimeZone=India"
	// gorm.Open(postgres.Open(dsn), &gorm.Config{}

}
