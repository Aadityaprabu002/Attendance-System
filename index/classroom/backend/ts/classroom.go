package classroom

import (
	connections "attsys/connections"
	key "attsys/env"
	"attsys/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

type ClassRoomDetails struct {
	ClassroomId  string
	DepartmentId string
	CourseId     string
}

type PageData struct {
	Todos []models.Classroom
}

var store = sessions.NewCookieStore([]byte(key.GetSecretKey()))

func CreateOrAppendClassRoom(TeacherId string, Classroom models.Classroom) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic("failed to establish connection with sql")
	}
	defer db.Close()
	query := fmt.Sprintf(`insert into classrooms (teacher_id,department_id,course_id)
	values('%s','%s','%s');`, TeacherId, Classroom.DepartmentId, Classroom.CourseId)

	query = strings.TrimSpace(query)
	fmt.Println(query)
	_, err = db.Exec(query)
	if err != nil {
		fmt.Println("Error creating a classroom")
		fmt.Println(err)
		return false
	}
	return true
}

func ClassroomSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	session, err := store.Get(r, "teacher")
	if err != nil {
		fmt.Println(err)
	}

	if !session.IsNew {
		if r.Method == "GET" {
			fmt.Println("GET")
			tmp, _ := template.ParseFiles("classroom/frontend/ts/session.html")
			r.ParseForm()
			ClassroomId := r.FormValue("ClassroomId")
			fmt.Println(ClassroomId)
			tmp.Execute(w, nil)
		}
	} else {
		http.Redirect(w, r, "/teacher/signin", http.StatusSeeOther)
	}

}
func Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	session, err := store.Get(r, "teacher")
	if err != nil {
		fmt.Println(err)
	}
	msg := models.Htmlresponse{
		Response: "",
		Status:   0,
	}
	if !session.IsNew {
		if r.Method == "GET" {
			fmt.Println("GET")
			tmp, _ := template.ParseFiles("classroom/frontend/ts/dashboard.html")
			data := PageData{
				Todos: GetClassrooms(session.Values["TEACHER_ID"].(string)),
			}
			fmt.Println("list of classrooms:")
			fmt.Print(data)
			tmp.Execute(w, data)

		} else if r.Method == "POST" {
			fmt.Println("POST")
			var params models.Classroom
			err := json.NewDecoder(r.Body).Decode(&params)
			if err != nil {
				fmt.Println(err)
			}
			if CreateOrAppendClassRoom(session.Values["TEACHER_ID"].(string), params) {
				msg.Status = 1
			} else {
				msg.Response = "Failed creating a classroom"
			}
			json.NewEncoder(w).Encode(msg)
		}
		return
	} else {
		http.Redirect(w, r, "/teacher/signin", http.StatusSeeOther)
	}

}
