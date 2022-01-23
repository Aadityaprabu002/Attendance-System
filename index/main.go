package main

import (
	adminas "attsys/admin/backend/as"
	admin_dept_course "attsys/admin/backend/department_and_courses"
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

func initRouter() {

	publicfs := http.FileServer(http.Dir("./static/public"))
	adminfs := http.FileServer(http.Dir("./static/admin"))
	admin := mux.NewRouter()
	admin.HandleFunc("/", adminas.Index)
	admin.HandleFunc("/admin/student/signup", adminss.Signup)
	admin.HandleFunc("/admin/student/handlestudent", adminss.HandleStudent).Methods("POST")
	admin.HandleFunc("/admin/teacher/signup", admints.TeacherSignup)
	admin.HandleFunc("/admin/teacher/handleteacher", admints.HandleTeacher).Methods("POST")
	admin.HandleFunc("/admin/department_and_courses", admin_dept_course.DepartmentAndCourse).Methods("GET")
	admin.HandleFunc("/admin/handledepartment", admin_dept_course.HandleDepartment).Methods("POST")
	admin.HandleFunc("/admin/handlecourse", admin_dept_course.HandleCourse).Methods("POST")

	admin.PathPrefix("/static/").Handler(http.StripPrefix("/static/", adminfs))

	public := mux.NewRouter()
	public.PathPrefix("/static/").Handler(http.StripPrefix("/static/", publicfs))
	public.HandleFunc("/", home.Homepage)
	public.HandleFunc("/student/signin", student.Signin)
	public.HandleFunc("/student/signin/complete_registration", student.CompleteRegistration)
	public.HandleFunc("/teacher/signin", teacher.Signin)
	public.HandleFunc("/teacher/signin/complete_registration", teacher.CompleteRegistration)
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
