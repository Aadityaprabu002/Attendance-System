package models

type Signup struct {
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
	Email     string   `json:"email"`
	Regnumber string   `json:"regnumber"`
	Password  []string `json:"password"`
	Image     string   `json:"image64"`
}

type Signin struct {
	Email     string `json:"email"`
	Regnumber string `json:"regnumber"`
	Password  string `json:"password"`
}

type User struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     string `json:"email"`
	Regnumber string `json:"regnumber"`
	Password  string `json:"-"`
}
