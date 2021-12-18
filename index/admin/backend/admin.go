package admin

import (
	models "attsys/models"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/lib/pq"

	"golang.org/x/crypto/bcrypt"
)

func TeacherSignup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Method == "GET" {
		fmt.Println("GET")
		tmp, _ := template.ParseFiles("admin/frontend/teacher/signup/signup.html")
		tmp.Execute(w, nil)
		return
	} else {
		fmt.Println("POST")
		msg := models.Htmlresponse{
			Response: "",
			Status:   0,
		}
		var params models.TeacherSignup
		err := json.NewDecoder(r.Body).Decode(&params)
		if err == nil {
			if params.Password[0] == params.Password[1] {
				encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(params.Password[0]), 8)
				newTeacher := models.Teacher{
					Email:        params.Email,
					Firstname:    params.Firstname,
					Lastname:     params.Lastname,
					TeacherId:    params.TeacherId,
					DepartmentId: params.DepartmentId,
					CourseId:     params.CourseId,
					Password:     string(encryptedPassword),
				}
				if !IsTeacherExist(newTeacher) {
					insertTeacher(newTeacher)
					msg.Status = 1
				} else {
					msg.Response = "Email or Teacher id  already exist!"
				}
			} else {
				msg.Response = "Password not matching!"
			}
		}
		json.NewEncoder(w).Encode(msg)

	}

}
