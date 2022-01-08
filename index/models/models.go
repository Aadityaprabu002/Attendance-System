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
	Regnumber   string    `json:"regnumber"`
	ClassroomId string    `json:"classroomid"`
	JoiningTime time.Time `json:"joiningtime"`
}

type SessionDetails struct {
	ClassroomId int       `json:"classroomid"`
	Start_time  time.Time `json:"start_time"`
	End_time    time.Time `json:"end_time"`
}
