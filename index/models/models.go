package models

import "time"

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
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Email        string `json:"email"`
	TeacherId    string `json:"teacherid"`
	DepartmentId string `json:"departmentid"`
	CourseId     string `json:"courseid"`
	Password     string `json:"-"`
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
	TeacherId    string `json:"teacherid"`
	ClassroomId  string `json:"classroomid"`
	DepartmentId string `json:"departmentid"`
	CourseId     string `json:"courseid"`
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
}

type PrettySession struct {
	SessionId  int
	Date       string
	Start_time string
	End_time   string
	Status     string
	SessionKey string
}

type StudentSessionDashBoard struct {
	SessionDetails PrettySession
	TeacherName    string
	DepartmentName string
	CourseName     string
}

type Attendance struct {
	PrettyTime    string
	ImageFilePath string
}

type AttendanceDetails struct {
	StudentName string
	Regnumber   string
	Attendance1 Attendance
	Attendance2 Attendance
	Attendance3 Attendance
}

type TeacherSessionDashBoard struct {
	SessionDetails PrettySession
	TeacherName    string
	DepartmentName string
	CourseName     string
	Attendees      []AttendanceDetails
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
