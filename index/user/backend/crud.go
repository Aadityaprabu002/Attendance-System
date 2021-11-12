package users

import (
	connections "attsys/connections"
	"attsys/models"
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"image/png"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// func reviveProcess() {
// 	if r := recover(); r != nil {
// 		fmt.Println("Process Revived!!")
// 	}
// }

func insertUser(newUser models.User) {
	// defer reviveProcess()
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	fmt.Println("Connection: " + conn)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
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
	// defer reviveProcess()
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	fmt.Println(conn)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
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
	// defer reviveProcess()
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

func saveUserImage(b64 string) {
	list := strings.Split(b64, ",")
	b := list[1]
	unbased, err := base64.StdEncoding.DecodeString(b)
	if err != nil {
		panic("Cannot decode base64")
	}
	r := bytes.NewReader(unbased)
	im, err := png.Decode(r)
	if err != nil {
		panic("Bad png")
	}
	f, err := os.OpenFile("../dummy/example.png", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		panic("Cannot open file")
	}
	err = png.Encode(f, im)
	if err != nil {
		panic("Error saving the file")
	}
}
