package admin

import (
	"attsys/connections"
	"attsys/models"
	"database/sql"
	"fmt"
)

func InsertDepartment(params models.Department) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	query := fmt.Sprintf(`insert into departments(department_id,department_name)
	values
	('%s','%s')`, params.DepartmentId, params.DepartmentName)
	_, err = db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}

func DeleteDepartment(DepartmentId string) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	query := fmt.Sprintf(`delete from departments where department_id = '%s'`, DepartmentId)
	_, err = db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}
func InsertCourse(params models.Course) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	query := fmt.Sprintf(`insert into courses (course_id,course_name)
	values ('%s','%s')`, params.CourseId, params.CourseName)

	_, err = db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}

func DeleteCourse(CourseId string) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	query := fmt.Sprintf(`delete from courses where course_id = '%s'`, CourseId)
	_, err = db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}

func IsCourseExist(CourseId string) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return true
	}
	defer db.Close()
	query := fmt.Sprintf(`select exists ( select 1  from courses where course_id = '%s')`, CourseId)
	result, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return true
	}
	var exist bool
	for result.Next() {
		result.Scan(&exist)
	}
	return exist
}

func IsDepartmentExist(DepartmentId string) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
		return true
	}
	defer db.Close()
	query := fmt.Sprintf(`select exists ( select 1  from departments where department_id = '%s')`, DepartmentId)
	result, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return true
	}
	var exist bool
	for result.Next() {
		result.Scan(&exist)
	}
	return exist
}
