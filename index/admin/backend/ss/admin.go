package admin

import (
	"attsys/models"
	"encoding/json"
	"html/template"
	"net/http"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Method == "GET" {
		tmp, _ := template.ParseFiles("admin/frontend/student/signup/signup.html")
		tmp.Execute(w, nil)
		return
	} else {
		msg := models.Htmlresponse{
			Response: "",
			Status:   0,
		}
		var params models.StudentSignup
		err := json.NewDecoder(r.Body).Decode(&params)
		if err == nil {
			if params.Password[0] == params.Password[1] {
				encryptedPassword := EncryptPassword(params.Password[0])
				newStudent := models.Student{
					// Email:     params.Email,
					Firstname: params.Firstname,
					Lastname:  params.Lastname,
					Regnumber: params.Regnumber,
					Password:  encryptedPassword,
				}
				if !IsStudentExist(newStudent.Regnumber) {
					// if ImagePath := SaveStudentImageData(newStudent.Regnumber, params.Image); ImagePath != "" {
					if InsertStudent(newStudent) {
						msg.Response = "Account created successfully!"
						msg.Status = 1
					} else {
						msg.Response = "Account creation failed!"
					}

				} else {
					msg.Response = "Student already Registered!"
				}
			} else {
				msg.Response = "Password not matching!"

			}
		}
		json.NewEncoder(w).Encode(msg)
	}

}
