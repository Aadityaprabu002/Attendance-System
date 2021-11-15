package users

import (
	models "attsys/models"
	"encoding/json"
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
	// defer reviveProcess()
	if r.Method == "GET" {
		fmt.Print("GET")
		tmp, _ := template.ParseFiles("user/frontend/signin/signin.html")
		tmp.Execute(w, nil)
		return
	} else {
		fmt.Print("POST")
		var params models.Signin
		err := json.NewDecoder(r.Body).Decode(&params)
		if err == nil {
			params := r.Form
			fmt.Println("Password:" + params["password"][0])
			user := models.User{
				Email:     params["email"][0],
				Password:  params["password"][0],
				Regnumber: params["regnumber"][0],
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
	// defer reviveProcess()
	msg := htmlresponse{
		Response: "",
	}

	if r.Method == "GET" {
		fmt.Print("GET")
		tmp, _ := template.ParseFiles("user/frontend/signup/signup.html")
		tmp.Execute(w, msg)
		return
	} else {
		fmt.Print("POST")
		var params models.Signup
		err := json.NewDecoder(r.Body).Decode(&params)
		if err == nil {
			if params.Password[0] == params.Password[1] {
				encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(params.Password[0]), 8)
				newUser := models.User{
					Email:     params.Email,
					Firstname: params.Firstname,
					Lastname:  params.Lastname,
					Regnumber: params.Regnumber,
					Password:  string(encryptedPassword),
				}
				if !checkUser(newUser.Email) {
					if saveUserImageData(newUser.Regnumber, params.Image) {
						insertUser(newUser)
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
		tmp, _ := template.ParseFiles("user/frontend/signup/matchface.html")
		tmp.Execute(w, nil)
		return
	}
}
