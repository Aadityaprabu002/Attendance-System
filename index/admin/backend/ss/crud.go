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

func InsertStudent(newUser models.Student) bool {

	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return false
	}
	defer db.Close()

	query := fmt.Sprintf(`
				Insert into students (firstname,lastname,password,regnumber,status)
				values('%s','%s','%s','%s',%d) 
				`, newUser.Firstname, newUser.Lastname, newUser.Password, newUser.Regnumber, 1)

	_, err = db.Exec(query)

	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to execute query")
		return false
	}
	return true
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
func UpdateStudentEmail(Regnumber string, newEmail string) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	query := fmt.Sprintf(`update students set email = '%s' where regnumber = '%s'`, newEmail, Regnumber)
	_, err = db.Query(query)
	return err == nil
}
func UpdateStudentPassword(Regnumber string, newPassword string) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	query := fmt.Sprintf(`update students set password = '%s' where regnumber = '%s'`, newPassword, Regnumber)
	_, err = db.Query(query)
	return err == nil

}
func UpdateStudentImagePath(Regnumber string, newImagePath string) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	query := fmt.Sprintf(`update students set image = '%s' where regnumber = '%s'`, newImagePath, Regnumber)
	_, err = db.Query(query)
	return err == nil

}
func UpdateStudentStatus(Regnumber string, status int) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	query := fmt.Sprintf(`update students set status = %d where regnumber = '%s'`, status, Regnumber)
	_, err = db.Query(query)
	return err == nil
}
func GetStudentAccountStatus(Regnumber string) int {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer db.Close()
	query := fmt.Sprintf(`select status from students where regnumber = '%s'`, Regnumber)
	var status int
	result, err := db.Query(query)
	if err != nil {
		return 0
	}
	for result.Next() {
		result.Scan(&status)
	}
	return status
}

func GetStudentDetails(Regnumber string) models.Student {
	var student models.Student
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	query := fmt.Sprintf(`select firstname,lastname,email,picture from students where regnumber = '%s'`, Regnumber)
	result, _ := db.Query(query)
	for result.Next() {
		result.Scan(&student.Firstname, &student.Lastname, &student.Email, &student.Image)
	}
	return student
}

func EncryptPassword(password string) string {
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(encryptedPassword)
}
