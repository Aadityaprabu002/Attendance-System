package student

import (
	admin "attsys/admin/backend/ss"
	key "attsys/env"
	models "attsys/models"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

var store = sessions.NewCookieStore([]byte(key.GetSecretKey()))

func IsLogged(r *http.Request) bool {
	session, _ := store.Get(r, "student")
	return !session.IsNew && session.Values["REG_NUMBER"] != nil
}

func Signin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Method == "GET" {
		if IsLogged(r) {
			http.Redirect(w, r, "/student/dashboard", http.StatusSeeOther)
		}
		fmt.Println("GET")
		tmp, _ := template.ParseFiles("student/frontend/signin/signin.html")
		tmp.Execute(w, nil)
		return
	} else {
		fmt.Println("POST")
		var params models.StudentSignin
		err := json.NewDecoder(r.Body).Decode(&params)
		msg := models.Htmlresponse{
			Response: "",
			Status:   0,
		}
		if err == nil {
			student := models.Student{
				Email:     params.Email,
				Password:  params.Password,
				Regnumber: params.Regnumber,
				Firstname: "",
				Lastname:  "",
			}
			if admin.IsStudentExist(student.Regnumber) {
				if admin.IsValidStudent(student) {
					fmt.Println("Valid Student!!")
					session, _ := store.Get(r, "student")
					if session.IsNew {
						session.Values["REG_NUMBER"] = student.Regnumber
						session.Options = &sessions.Options{
							Path:     "/",
							MaxAge:   3600 * 8, // 8 hours session timing
							HttpOnly: true,
						}
						err := session.Save(r, w)
						if err != nil {
							fmt.Println(err)
						}
					}
					msg.Status = 1
				} else {
					msg.Response = "Invalid username or password!!"

				}
			} else {
				msg.Response = "Invalid username or password!!"

			}
		}
		fmt.Println(msg)
		json.NewEncoder(w).Encode(msg)
	}
}

func Signout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		session, _ := store.Get(r, "student")
		if IsLogged(r) {
			if session.Values["REG_NUMBER"] != nil {
				session.Values["REG_NUMBER"] = nil
			}
			if session.Values["SESSION_KEY"] != nil {
				session.Values["SESSION_KEY"] = nil
			}
		}
		session.Options.MaxAge = -1
		session.Save(r, w)
		http.Redirect(w, r, "/student/signin", http.StatusSeeOther)
	}

}
