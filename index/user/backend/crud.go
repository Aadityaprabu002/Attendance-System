package users

import (
	connections "attsys/connections"
	"attsys/models"
	"database/sql"
	"fmt"
	"io/ioutil"
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
		panic("failed to establish connection with sql")
	}
	defer db.Close()

	query := fmt.Sprintf(`
				Insert into students (email,firstname,lastname,password,regnumber)
				values('%s','%s','%s','%s','%s') 
				`, newUser.Email, newUser.Firstname, newUser.Lastname, newUser.Password, newUser.Regnumber)

	query = strings.TrimSpace(query)
	fmt.Print(query)
	_, err = db.Exec(query)

	if err != nil {
		fmt.Println(err)
		panic("Failed to execute query")
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
	query := fmt.Sprintf(`select exists(select 1 from students where email = '%s')`, email)
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
	query := fmt.Sprintf(`select password from students where email = '%s' limit 1`, user.Email)
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

func saveUserImageData(regnumber string, b64 string) bool {

	dir := fmt.Sprintf("../database/%s/", regnumber)
	err := os.Mkdir(dir, 0777)
	if err != nil {
		fmt.Println("Failed creating directory for new student")
	}

	tempFile, err := ioutil.TempFile(dir, "*.txt")
	if err != nil {
		fmt.Println("Error creating user image file")
		return false
	}
	defer tempFile.Close()
	_, err = tempFile.Write([]byte(b64))
	if err != nil {
		fmt.Println("Error writing user image file")
		return false
	}
	return true
}
