package admin

import (
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Method == "GET" {
		tmp, _ := template.ParseFiles("admin/frontend/admin/index.html")
		tmp.Execute(w, nil)
	}
}
