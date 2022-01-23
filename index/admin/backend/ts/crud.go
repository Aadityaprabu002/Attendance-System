package admin

import (
	connections "attsys/connections"
	"attsys/models"
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func insertTeacher(newTeacher models.Teacher) {

	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	fmt.Println("Connection: " + conn)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic("failed to establish connection with sql")
	}
	defer db.Close()

	query := fmt.Sprintf(`
				Insert into teachers (firstname,lastname,password,teacher_id)
				values('%s','%s','%s','%s') 
				`, newTeacher.Firstname, newTeacher.Lastname, newTeacher.Password, newTeacher.TeacherId)

	query = strings.TrimSpace(query)
	fmt.Print(query)
	_, err = db.Exec(query)

	if err != nil {
		fmt.Println(err)
		panic("Failed to execute query")
	}

}

func IsTeacherExist(TeacherId string) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	fmt.Println(conn)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	query := fmt.Sprintf(`select exists(select 1 from teachers where teacher_id = '%s')`, TeacherId)

	result, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		return false
	}

	var teacheridExist bool
	for result.Next() {
		result.Scan(&teacheridExist)
	}

	return teacheridExist
}

func IsValidTeacher(user models.Teacher) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	query := fmt.Sprintf(`select password from teachers where teacher_id = '%s' limit 1`, user.TeacherId)
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
			fmt.Println(flag)
			return flag == nil
		}
	}
	return false

}

func GetAllTeacher() []models.TeachersDetails {
	var Teachers []models.TeachersDetails
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")

	}
	defer db.Close()
	query := `select teacher_id,concat(firstname,' ',lastname),email from teachers;`

	result, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return Teachers
	}
	for result.Next() {
		var temp models.TeachersDetails

		result.Scan(&temp.TeacherId, &temp.Teachername, &temp.Email)

		Teachers = append(Teachers, temp)
	}
	return Teachers
}

func RemoveTeacher(TeacherId string) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	query := fmt.Sprintf(`delete from teachers where teacher_id = '%s'`, TeacherId)
	_, err = db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}

func GetTeacherAccountStatus(TeacherId string) int {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	defer db.Close()
	query := fmt.Sprintf(`select status from teachers where teacher_id = '%s'`, TeacherId)
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

func GetTeacherDetails(TeacherId string) models.Teacher {
	var teacher models.Teacher
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return teacher
	}
	defer db.Close()

	query := fmt.Sprintf(`select firstname,lastname,email from teachers where teacher_id = '%s'`, TeacherId)
	result, _ := db.Query(query)
	for result.Next() {
		result.Scan(&teacher.Firstname, &teacher.Lastname, &teacher.Email)
	}
	return teacher
}

func UpdateTeacherEmail(TeacherId string, newEmail string) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	query := fmt.Sprintf(`update teachers set email = '%s' where teacher_id = '%s'`, newEmail, TeacherId)
	_, err = db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}
func UpdateTeacherPassword(TeacherId string, newPassword string) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	query := fmt.Sprintf(`update teachers set password = '%s' where teacher_id = '%s'`, newPassword, TeacherId)
	_, err = db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}

func UpdateTeacherStatus(TeacherId string, status int) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	query := fmt.Sprintf(`update teachers set status = %d where teacher_id = '%s'`, status, TeacherId)
	_, err = db.Query(query)

	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}
func EncryptPassword(password string) string {
	encryptedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(encryptedPassword)
}
