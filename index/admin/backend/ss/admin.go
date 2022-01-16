package admin

import (
	"attsys/models"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method == "GET" {
		fmt.Println("GET")
		tmp, _ := template.ParseFiles("admin/frontend/student/signup/signup.html")
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
					// Email:     params.Email,
					Firstname: params.Firstname,
					Lastname:  params.Lastname,
					Regnumber: params.Regnumber,
					Password:  string(encryptedPassword),
				}
				if !IsStudentExist(newStudent.Regnumber) {
					if ImagePath := SaveStudentImageData(newStudent.Regnumber, params.Image); ImagePath != "" {
						InsertStudent(newStudent)
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
