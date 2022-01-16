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
				Insert into teachers (email,firstname,lastname,password,teacher_id,department_id,course_id)
				values('%s','%s','%s','%s','%s','%s','%s') 
				`, newTeacher.Email, newTeacher.Firstname, newTeacher.Lastname, newTeacher.Password, newTeacher.TeacherId, newTeacher.DepartmentId, newTeacher.CourseId)

	query = strings.TrimSpace(query)
	fmt.Print(query)
	_, err = db.Exec(query)

	if err != nil {
		fmt.Println(err)
		panic("Failed to execute query")
	}

}

func IsTeacherExist(teacher models.Teacher) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	fmt.Println(conn)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer db.Close()
	query := fmt.Sprintf(`select exists(select 1 from teachers where teacher_id = '%s')`, teacher.TeacherId)
	query = strings.TrimSpace(query)

	result, err := db.Query(query)

	if err != nil {
		panic(err)
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
