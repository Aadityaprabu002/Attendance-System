package main

import (
	adminas "attsys/admin/backend/as"
	adminss "attsys/admin/backend/ss"
	admints "attsys/admin/backend/ts"
	student_classroom "attsys/classroom/backend/ss"
	teacher_classroom "attsys/classroom/backend/ts"
	home "attsys/home/backend"
	student "attsys/student/backend"
	teacher "attsys/teacher/backend"
	"fmt"
	"log"
	"net/http"
	"time"

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

	publicfs := http.FileServer(http.Dir("./static/public"))
	adminfs := http.FileServer(http.Dir("./static/admin"))
	admin := mux.NewRouter()
	admin.HandleFunc("/", adminas.Index)
	admin.HandleFunc("/admin/student/signup", adminss.Signup)
	admin.HandleFunc("/admin/teacher/signup", admints.TeacherSignup)
	admin.PathPrefix("/static/").Handler(http.StripPrefix("/static/", adminfs))

	public := mux.NewRouter()
	public.PathPrefix("/static/").Handler(http.StripPrefix("/static/", publicfs))
	public.HandleFunc("/", home.Homepage)
	public.HandleFunc("/student/signin", student.Signin)
	public.HandleFunc("/student/signin/complete_registration", student.CompleteRegistration)
	public.HandleFunc("/teacher/signin", teacher.Signin)
	public.HandleFunc("/teacher/signout", teacher.Signout)
	public.HandleFunc("/teacher/dashboard", teacher_classroom.Dashboard)
	public.HandleFunc("/teacher/dashboard/classroomdashboard/{ClassroomId}", teacher_classroom.ClassroomDashboard)
	public.HandleFunc("/teacher/dashboard/classroomdashboard/{ClassroomId}/handlestudents", teacher_classroom.Handlestudents).Methods("POST")
	public.HandleFunc("/teacher/dashboard/classroomdashboard/{ClassroomId}/sessiondashboard/{SessionId}", teacher_classroom.SessionDashboard)
	public.HandleFunc("/teacher/dashboard/classroomdashboard/{ClassroomId}/sessiondashboard/{SessionId}/postattendance", teacher_classroom.PostAttendance)
	public.HandleFunc("/student/dashboard", student_classroom.Dashboard)
	public.HandleFunc("/student/dashboard/session", student_classroom.SessionDashboard).Methods("GET")
	public.HandleFunc("/student/dashboard/session/timerdetails", student_classroom.SessionDetails).Methods("GET")
	public.HandleFunc("/student/dashboard/session/postattendance", student_classroom.PostAttendance).Methods("POST")
	public.HandleFunc("/student/dashboard/session/endsession", student_classroom.EndSession).Methods("GET")
	public.HandleFunc("/student/signout", student.Signout)
	go func() {
		log.Fatal(http.ListenAndServe(":4040", admin))
	}()
	log.Fatal(http.ListenAndServe(":8080", public))

}

func main() {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	now := time.Now().In(loc)
	fmt.Println("Location : ", loc, " Time : ", now)

	initRouter()
}
