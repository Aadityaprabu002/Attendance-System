package users

import (
	connections "attsys/connections"
	"attsys/models"
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func reviveProcess() {
	if r := recover(); r != nil {
		fmt.Println("Process Revived!!")
	}
}

func insertUser(newUser models.User) {
	defer reviveProcess()
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	fmt.Println("Connection: " + conn)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	query := fmt.Sprintf(`
				Insert into users (email,firstname,lastname,password)
				values('%s','%s','%s','%s')
				`, newUser.Email, newUser.Firstname, newUser.Lastname, newUser.Password)

	query = strings.TrimSpace(query)
	_, err = db.Exec(query)

	if err != nil {
		panic(err)
	}

}

func checkUser(email string) bool {
	defer reviveProcess()
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	query := fmt.Sprintf(`select exists(select 1 from users where email = '%s')`, email)
	query = strings.TrimSpace(query)

	result, err := db.Query(query)

	if err != nil {
		panic(err)
	}

	var exist bool
	for result.Next() {
		result.Scan(&exist)
	}

	return exist
}

func validUser(user models.User) bool {
	defer reviveProcess()
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	query := fmt.Sprintf(`select password from users where email = '%s' limit 1`, user.Email)
	query = strings.TrimSpace(query)
	fmt.Println(query)
	result, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	var encryptedPassword []byte

	for result.Next() {
		err := result.Scan(&encryptedPassword)
		if err != nil {
			panic(err)
		} else {
			flag := bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(user.Password))
			return flag == nil
		}
	}
	return false

}
