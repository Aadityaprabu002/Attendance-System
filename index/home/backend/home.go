package home

import (
	"html/template"
	"net/http"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	tmp, _ := template.ParseFiles("home/frontend/home.html")
	tmp.Execute(w, nil)
}
