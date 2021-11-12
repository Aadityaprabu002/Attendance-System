package users

import (
	models "attsys/models"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type htmlresponse struct {
	Response string
}

func Signin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	defer reviveProcess()
	if r.Method == "GET" {
		fmt.Print("GET")
		tmp, _ := template.ParseFiles("user/frontend/signin.html")
		tmp.Execute(w, nil)
		return
	} else {
		fmt.Print("POST")
		err := r.ParseForm()
		if err == nil {
			params := r.Form
			fmt.Println("Password:" + params["password"][0])
			user := models.User{
				Email:     params["email"][0],
				Password:  params["password"][0],
				Firstname: "",
				Lastname:  "",
			}
			if checkUser(user.Email) {
				if validUser(user) {
					fmt.Println("Valid user!!")
				} else {
					fmt.Println("Invalid username or password!!")
				}
			} else {
				fmt.Println("Invalid username or password!!")
			}
		}

		tmp, _ := template.ParseFiles("user/frontend/signin.html")
		tmp.Execute(w, nil)
	}
}
func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	defer reviveProcess()
	msg := htmlresponse{
		Response: "",
	}

	if r.Method == "GET" {
		fmt.Print("GET")
		tmp, _ := template.ParseFiles("user/frontend/signup.html")
		tmp.Execute(w, msg)
		return
	} else {
		fmt.Print("POST")
		err := r.ParseForm()

		if err == nil {
			params := r.Form
			fmt.Println(params["image"][0])
			fmt.Println("Password:" + params["password"][0])
			if params["password"][0] == params["password"][1] {
				encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(params["password"][0]), 8)
				newUser := models.User{
					Email:     params["email"][0],
					Firstname: params["fname"][0],
					Lastname:  params["lname"][0],
					Password:  string(encryptedPassword),
				}
				if !checkUser(newUser.Email) {
					insertUser(newUser)
				} else {
					msg.Response = "Email already exist!"
				}
			} else {
				msg.Response = "Password not matching!"
			}

		}
		tmp, _ := template.ParseFiles("user/frontend/signup.html")

		tmp.Execute(w, nil)
	}

}
