package teacher

import (
	admin "attsys/admin/backend/ts"
	key "attsys/env"
	"attsys/models"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

func IsLogged(r *http.Request) bool {
	session, _ := store.Get(r, "teacher")
	return !session.IsNew && session.Values["TEACHER_ID"] != nil
}

var store = sessions.NewCookieStore([]byte(key.GetSecretKey()))

func Signin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if IsLogged(r) {
		http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
	} else {

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
				if admin.IsTeacherExist(teacher.TeacherId) {
					if admin.IsValidTeacher(teacher) {
						fmt.Println("Valid teacher!!")
						session, _ := store.Get(r, "teacher")
						if session.IsNew {
							session.Values["TEACHER_ID"] = teacher.TeacherId
							session.Options = &sessions.Options{
								Path:     "/",
								MaxAge:   3600 * 4, // 4 hours session timing
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
}

func CompleteRegistration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if IsLogged(r) {
		session, _ := store.Get(r, "teacher")
		TeacherId := session.Values["TEACHER_ID"].(string)
		AccountStatus := admin.GetTeacherAccountStatus(TeacherId)
		if AccountStatus == 1 {
			if r.Method == "GET" {
				tmp, _ := template.ParseFiles("teacher/frontend/signin/completeRegistration.html")
				res := admin.GetTeacherDetails(TeacherId)
				tmp.Execute(w, res)
			} else if r.Method == "POST" {
				res := models.Htmlresponse{
					Response: "Failed Registering",
					Status:   0,
				}

				params := models.StudentSignup{}
				err := json.NewDecoder(r.Body).Decode(&params)
				fmt.Println(err)

				if err == nil {
					if params.Password[0] == params.Password[1] {
						encryptedPassword := admin.EncryptPassword(params.Password[0])
						if admin.UpdateTeacherEmail(TeacherId, params.Email) && admin.UpdateTeacherPassword(TeacherId, encryptedPassword) {
							if admin.UpdateTeacherStatus(TeacherId, 2) {
								res.Response = "Registration completed successfully"
								res.Status = 1
							}
						}

					}
				}
				json.NewEncoder(w).Encode(res)
			}
		} else if AccountStatus == 2 {
			http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/teacher/signout", http.StatusSeeOther)
		}

	}
}

func Signout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "GET" {
		session, _ := store.Get(r, "teacher")
		if IsLogged(r) {
			if session.Values["TEACHER_ID"] != nil {
				session.Values["REG_NUMBER"] = nil
			}
		}
		session.Options.MaxAge = -1
		session.Save(r, w)
		http.Redirect(w, r, "/teacher/signin", http.StatusSeeOther)
	}

}
