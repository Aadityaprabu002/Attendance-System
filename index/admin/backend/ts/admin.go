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
		res := struct {
			Teachers []models.TeachersDetails
		}{}
		res.Teachers = GetAllTeacher()
		tmp.Execute(w, res)

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
					Firstname: params.Firstname,
					Lastname:  params.Lastname,
					TeacherId: params.TeacherId,
					Password:  string(encryptedPassword),
				}
				if !IsTeacherExist(newTeacher.TeacherId) {
					insertTeacher(newTeacher)
					msg.Status = 1
					msg.Response = "Successfully account created"
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
func HandleTeacher(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		msg := models.Htmlresponse{
			Response: "Failed to remove Teacher",
			Status:   0,
		}
		params := struct {
			TeacherId string `json:"teacherid"`
		}{}
		json.NewDecoder(r.Body).Decode(&params)
		if IsTeacherExist(params.TeacherId) {
			if RemoveTeacher(params.TeacherId) {
				msg.Response = "Successfully removed Teacher !"
				msg.Status = 1
			}
		} else {
			msg.Response = "Teacher does not exists"
		}

		json.NewEncoder(w).Encode(msg)
	}
}
