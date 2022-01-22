package models

import (
	"html/template"
	"time"
)

type Htmlresponse struct {
	Response string
	Status   int
}

// student models
type Student struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Regnumber string `json:"regnumber"`
	Image     string `json:"image64"`
	Password  string `json:"-"`
}

type StudentSignup struct {
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
	Email     string   `json:"email"`
	Regnumber string   `json:"regnumber"`
	Password  []string `json:"password"`
	Image     string   `json:"image64"`
}

type StudentSignin struct {
	Email     string `json:"email"`
	Regnumber string `json:"regnumber"`
	Password  string `json:"password"`
}

// teacher models
type Teacher struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	TeacherId string `json:"teacherid"`
	Password  string `json:"-"`
}

type TeacherSignup struct {
	Firstname    string   `json:"firstname"`
	Lastname     string   `json:"lastname"`
	Email        string   `json:"email"`
	TeacherId    string   `json:"teacherid"`
	DepartmentId string   `json:"departmentid"`
	CourseId     string   `json:"courseid"`
	Password     []string `json:"password"`
}
type TeacherSignin struct {
	TeacherId string `json:"teacherid"`
	Password  string `json:"password"`
}

// classroom models
type Classroom struct {
	TeacherId    string    `json:"teacherid"`
	ClassroomId  string    `json:"classroomid"`
	DepartmentId string    `json:"departmentid"`
	CourseId     string    `json:"courseid"`
	From         time.Time `json:"from"`
	To           time.Time `json:"to"`
}

//  joinee models
type Joinee struct {
	SessionKey string `json:"session_key"`
}

type Session struct {
	ClassroomId int       `json:"classroomid"`
	SessionId   int       `json:"sessionid"`
	Date        time.Time `json:"date"`
	Start_time  time.Time `json:"start_time"`
	End_time    time.Time `json:"end_time"`
	Status      string    `json:"status"`
	SessionKey  string
	Reviewed    bool
}

// type PrettySession struct {
// 	SessionId  int
// 	Date       string
// 	Start_time string
// 	End_time   string
// 	Status     string
// 	SessionKey string
// 	Reviewed   bool
// }

type SessionDashBoardDetails struct {
	SessionDetails Session
	TeacherName    string
	DepartmentName string
	CourseName     string
}

type Attendance struct {
	Time          time.Time
	ImageFilePath template.URL
}

type AttendanceDetails struct {
	StudentName string
	Regnumber   string
	FairImage   template.URL
	Attendance1 Attendance
	Attendance2 Attendance
	Attendance3 Attendance
	IsPresent   bool
}

type TeacherSessionDashBoardDetails struct {
	SessionDb SessionDashBoardDetails
	Attendees []AttendanceDetails
}

type StudentSessionTimerDetails struct {
	StartTime time.Time
	EndTime   time.Time
	PopUp1    time.Time
	PopUp2    time.Time
	PopUp3    time.Time
}

type PostAttendance struct {
	AttNum int       `json:"attendance_number"`
	Time   time.Time `json:"attendance_time"`
	Image  string    `json:"image64"`
}

type ReviewAttendance struct {
	Regnumber string `json:"regnumber"`
	IsPresent bool   `json:"is_present"`
}

type Course struct {
	CourseId   string
	CourseName string
}
type Department struct {
	DepartmentId   string
	DepartmentName string
}
type ClassroomTableData struct {
	Classrooms []Classroom
	Courses    []Course
	Department []Department
}

type StudentsDetails struct {
	Studentname string
	Regnumber   string
	Image       template.URL
	Email       string
}
type ClassroomTableDetails struct {
	Sessions []Session
	Students []StudentsDetails
}
