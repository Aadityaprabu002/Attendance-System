package users

import (
	models "attsys/models"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"golang.org/x/crypto/bcrypt"
)

type htmlresponse struct {
	Response string
}

func GETSECRETKEY() string {

	err := godotenv.Load("./env/safe.env")
	if err != nil {
		fmt.Println("Error Loading env file")
		fmt.Println(err)
	}
	return os.Getenv("SECRET_KEY")
}

var store = sessions.NewCookieStore([]byte(GETSECRETKEY()))

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
		msg := htmlresponse{
			Response: "",
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
							MaxAge:   10, // 9 hours session timing
							HttpOnly: true,
						}
						err := session.Save(r, w)
						if err != nil {
							fmt.Println(err)
						}
					}
					fmt.Println(session.Values["REG_NUMBER"])
				} else {
					msg.Response = "Invalid username or password!!"
				}
			} else {
				msg.Response = "Invalid username or password!!"
			}
		}
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
		msg := htmlresponse{
			Response: "",
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

func MatchFace(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Print("GET")
		tmp, _ := template.ParseFiles("student/frontend/signup/matchface.html")
		tmp.Execute(w, nil)
		return
	}
}
