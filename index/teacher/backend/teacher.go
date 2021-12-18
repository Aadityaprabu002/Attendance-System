package teacher

import (
	admin "attsys/admin/backend"
	key "attsys/env"
	"attsys/models"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(key.GetSecretKey()))

func Signin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method == "GET" {
		fmt.Println("GET")
		tmp, _ := template.ParseFiles("teacher/frontend/signin/signin.html")
		tmp.Execute(w, nil)
		return
	} else {
		fmt.Println("POST")
		var params models.TeacherSignin
		err := json.NewDecoder(r.Body).Decode(&params)
		msg := models.Htmlresponse{
			Response: "",
			Status:   0,
		}
		if err == nil {
			teacher := models.Teacher{
				Password:  params.Password,
				TeacherId: params.TeacherId,
				Email:     "",
				Firstname: "",
				Lastname:  "",
			}
			if admin.IsTeacherExist(teacher) {
				if admin.IsValidTeacher(teacher) {
					fmt.Println("Valid teacher!!")
					session, _ := store.Get(r, "teacher")
					if session.IsNew {
						session.Values["TEACHER_ID"] = teacher.TeacherId
						session.Options = &sessions.Options{
							Path:     "/",
							MaxAge:   3600, // 1 hours session timing
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
				msg.Response = "Account does not exist!!"
			}
		}
		fmt.Println(msg)
		json.NewEncoder(w).Encode(msg)
	}
}
