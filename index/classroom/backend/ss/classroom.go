package classroom

import (
	key "attsys/env"
	models "attsys/models"
	student "attsys/student/backend"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(key.GetSecretKey()))

func IsClassroomSet(r *http.Request) bool {
	session, _ := store.Get(r, "student")
	return session.Values["CLASSROOM_ID"] != nil && session.Values["JOINED_AT"] != nil
}

func JoinClassroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if !student.IsLogged(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		fmt.Println("GET")
		tmp, _ := template.ParseFiles("classroom/frontend/ss/join.html")
		tmp.Execute(w, nil)
		return
	} else {
		fmt.Println("POST")
		session, _ := store.Get(r, "student")

		params := models.Joinee{
			Regnumber: session.Values["REG_NUMBER"].(string),
		}

		err := json.NewDecoder(r.Body).Decode(&params)
		fmt.Println("params:")
		fmt.Println(params)
		msg := models.Htmlresponse{
			Response: "",
			Status:   0,
		}
		if err == nil {
			if IsClassRoomExist(params.ClassroomId) {
				session, _ := store.Get(r, "student")
				if !session.IsNew {
					session.Values["CLASSROOM_ID"] = params.ClassroomId
					session.Values["JOINED_AT"] = params.JoiningTime.String()
					err := session.Save(r, w)
					if err == nil {
						msg.Status = 1
						fmt.Println("class exist")
					} else {
						fmt.Println(err)
						msg.Response = "error joining classroom!"
					}
				}
			} else {
				msg.Response = "classroom does not exist"
			}
		} else {
			msg.Response = "Faulty input !"
		}

		fmt.Println("msg:")
		fmt.Println(msg)
		json.NewEncoder(w).Encode(msg)
		return
	}
}
func LoadClassroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if !student.IsLogged(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if !IsClassroomSet(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		fmt.Println("GET")
		tmp, _ := template.ParseFiles("classroom/frontend/ss/join.html")
		tmp.Execute(w, nil)
	} else {
		fmt.Println("POST")

	}
}
func LeaveClassroom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if !student.IsLogged(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if !IsClassroomSet(r) {
		http.Redirect(w, r, "/signin", http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		session, _ := store.Get(r, "student")
		if !session.IsNew {
			session.Values["CLASSROOM_ID"] = nil
			session.Values["JOINED_AT"] = nil
			err := session.Save(r, w)
			if err != nil {
				fmt.Println(err)
			}
		}
		return
	}
}
