package admin

import (
	"attsys/connections"
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

func InsertStudent(newUser models.Student) {
	// defer reviveProcess()
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
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

func IsStudentExist(Regnumber string) bool {

	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	fmt.Println(conn)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error connecting to database!")
	}
	defer db.Close()
	query := fmt.Sprintf(`select exists(select 1 from students where regnumber = '%s')`, Regnumber)
	query = strings.TrimSpace(query)

	result, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Error in checking student already exists")
	}

	var exist bool
	for result.Next() {
		result.Scan(&exist)
	}

	return exist
}

func IsValidStudent(user models.Student) bool {

	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	query := fmt.Sprintf(`select password from students where regnumber = '%s' limit 1`, user.Regnumber)
	query = strings.TrimSpace(query)
	fmt.Println(query)
	result, err := db.Query(query)

	if err != nil {
		panic(err)
	}
	var encryptedPassword string

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

func SaveStudentImageData(regnumber string, ImageData string) string {
	fpath := fmt.Sprintf("../database/students/%s/", regnumber)
	err := os.Mkdir(fpath, 0777)
	if err != nil {
		fmt.Println("Failed creating directory for new student")
	}

	ImageData = strings.Split(ImageData, ",")[1]
	DecodedImageData, err := base64.StdEncoding.DecodeString(ImageData)
	if err != nil {
		fmt.Println("Bad base 64 string!")
		return ""
	}
	ImgReader := bytes.NewReader(DecodedImageData)
	Image, err := png.Decode(ImgReader)
	if err != nil {
		fmt.Println("Bad image")
		return ""
	}
	ImagePath := fmt.Sprintf("%simage.png", fpath)
	ImageFile, err := os.OpenFile(ImagePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("Cannot open file")
		return ""
	}
	png.Encode(ImageFile, Image)
	return ImagePath
}
