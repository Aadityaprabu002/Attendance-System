package models

type Signup struct {
	Firstname string   `json:"fname"`
	Lastname  string   `json:"lname"`
	Email     string   `json:"email"`
	Password  []string `json:"password"`
	Image     string   `json:"image64"`
}

type Signin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"-"`
}
