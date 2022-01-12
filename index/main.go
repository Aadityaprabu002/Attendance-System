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
	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	initRouter()
}
