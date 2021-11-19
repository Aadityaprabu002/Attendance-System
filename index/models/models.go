package models

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

type Student struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Regnumber string `json:"regnumber"`
	Password  string `json:"-"`
}
