package main

import (
	admin "attsys/admin/backend"
	student_classroom "attsys/classroom/backend/ss"
	teacher_classroom "attsys/classroom/backend/ts"
	home "attsys/home/backend"
	student "attsys/student/backend"
	teacher "attsys/teacher/backend"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// var dummy = sessions.NewCookieStore([]byte(key.GetSecretKey()))

// func handleSession(w http.ResponseWriter, r *http.Request) {
// 	ssn, _ := dummy.Get(r, "dummy")

// 	if ssn.IsNew {
// 		ssn.Values["name"] = "value"
// 		ssn.Options = &sessions.Options{
// 			Domain:   "/",
// 			MaxAge:   10,
// 			HttpOnly: true,
// 		}
// 		ssn.Save(r, w)
// 		fmt.Println("New session! Existing Session Expired!")
// 		fmt.Println(ssn)
// 	} else {
// 		fmt.Println("")
// 		fmt.Println("Old session!")
// 		fmt.Println(ssn)
// 	}

// }
func initRouter() {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	r.HandleFunc("/", home.Homepage)
	r.HandleFunc("/student/signin", student.Signin)
	r.HandleFunc("/student/signup", student.Signup)
	r.HandleFunc("/admin/teacher/signup", admin.TeacherSignup)
	r.HandleFunc("/teacher/signin", teacher.Signin)
	r.HandleFunc("/teacher/dashboard", teacher_classroom.Dashboard)
	r.HandleFunc("/teacher/dashboard/classroomdashboard/{ClassroomId}", teacher_classroom.ClassroomDashboard)
	r.HandleFunc("/teacher/dashboard/classroomdashboard/{ClassroomId}/sessiondashboard/{SessionId}", teacher_classroom.SessionDashboard)
	r.HandleFunc("/student/dashboard", student_classroom.Dashboard)
	r.HandleFunc("/student/dashboard/session", student_classroom.SessionDashboard)
	r.HandleFunc("/student/dashboard/session/timerdetails", student_classroom.SessionDetails).Methods("GET")
	r.HandleFunc("/student/dashboard/session/postattendance", student_classroom.PostAttendance).Methods("GET")
	// r.HandleFunc("/", handleSession)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	initRouter()
}
