package student

import (
	key "attsys/env"
	models "attsys/models"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"

	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte(key.GetSecretKey()))

func IsLogged(r *http.Request) bool {
	session, _ := store.Get(r, "student")
	return !session.IsNew && session.Values["REG_NUMBER"] != nil
}
func Signin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// defer reviveProcess()
	if r.Method == "GET" {
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
			if isStudentExist(student) {
				if isValidStudent(student) {
					fmt.Println("Valid Student!!")
					session, _ := store.Get(r, "student")
					if session.IsNew {
						session.Values["REG_NUMBER"] = student.Regnumber
						session.Options = &sessions.Options{
							Path:     "/",
							MaxAge:   3600 * 9, // 9 hours session timing
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
func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method == "GET" {
		fmt.Println("GET")
		tmp, _ := template.ParseFiles("student/frontend/signup/signup.html")
		tmp.Execute(w, nil)
		return
	} else {
		fmt.Println("POST")
		msg := models.Htmlresponse{
			Response: "",
			Status:   0,
		}
		var params models.StudentSignup
		err := json.NewDecoder(r.Body).Decode(&params)
		if err == nil {
			if params.Password[0] == params.Password[1] {
				encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(params.Password[0]), 8)
				newStudent := models.Student{
					Email:     params.Email,
					Firstname: params.Firstname,
					Lastname:  params.Lastname,
					Regnumber: params.Regnumber,
					Password:  string(encryptedPassword),
				}
				if !isStudentExist(newStudent) {
					if saveStudentImageData(newStudent.Regnumber, params.Image) {
						insertStudent(newStudent)
						msg.Response = "Successfully account created!"
						msg.Status = 1
					} else {
						fmt.Println("Failed creating a record!")
					}

				} else {
					msg.Response = "Email already exist!"

				}
			} else {
				msg.Response = "Password not matching!"

			}
		}
		json.NewEncoder(w).Encode(msg)

	}

}

func Signout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "student")
	if err != nil {
		fmt.Println(err)
	}
	session.Options.MaxAge = -1
	session.Save(r, w)
}
func MatchFace(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Print("GET")
		tmp, _ := template.ParseFiles("student/frontend/signup/matchface.html")
		tmp.Execute(w, nil)
		session, err := store.Get(r, "student")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(session.Values["REG_NUMBER"])
		return
	}
}
