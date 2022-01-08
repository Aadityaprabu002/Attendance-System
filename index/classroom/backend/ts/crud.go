package classroom

import (
	connections "attsys/connections"
	"attsys/models"
	"database/sql"
	"fmt"
	"strings"
)

func IsClassRoomExist(ClassroomId string) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic("failed to establish connection with sql")
	}
	defer db.Close()
	query := fmt.Sprintf(`select exists(select 1 from classrooms where classroom_id = '%s')`, ClassroomId)
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

func GetClassrooms(TeacherId string) []models.Classroom {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
	}
	defer db.Close()
	query := fmt.Sprintf(`select classroom_id,department_id,course_id from classrooms
	where teacher_id = '%s'`, TeacherId)
	result, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	var temp models.Classroom
	var listOfClassrooms []models.Classroom
	for result.Next() {
		result.Scan(&temp.ClassroomId, &temp.DepartmentId, &temp.CourseId)
		listOfClassrooms = append(listOfClassrooms, temp)

	}
	return listOfClassrooms
}

func createUnqiueSession(newSession models.SessionDetails) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return false
	}

	defer db.Close()
	query := fmt.Sprintf("select * from classrooms where classroom_id = %d", newSession.ClassroomId)
	result, err := db.Query(query)
	if err != nil {
		fmt.Println("Error occurred when getting classroom details!")
		return false
	}
	var c models.Classroom
	for result.Next() {
		result.Scan(&c.TeacherId, &c.DepartmentId, &c.CourseId, &c.ClassroomId)
	}
	sessionUniqueCode := c.ClassroomId + "_" + c.DepartmentId + "_" + c.CourseId
	fmt.Println(sessionUniqueCode)
	query = fmt.Sprintf("insert into sessions(session_code,session_date,start_time,end_time,classroom_id) values ('%s','%s','%s','%s',%d)", sessionUniqueCode, newSession.End_time.Format("2006-01-02"), newSession.Start_time.Format("15:04:05"), newSession.End_time.Format("15:04:05"), newSession.ClassroomId)
	fmt.Println(query)
	_, err = db.Query(query)
	if err != nil {

		fmt.Println("Error creating a session")
		fmt.Println(err)
		return false
	}
	return true

}

func checkForSession(ClassroomId int) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
	}
	defer db.Close()

	query := `select now()::time < max(end_time) from sessions
where session_date = (
select max(session_date) from sessions) and classroom_id = %d`

	query = fmt.Sprintf(query, ClassroomId)
	query = strings.TrimSpace(query)

	result, err := db.Query(query)
	if err != nil {
		fmt.Println("Error occurred during checking if session already exists")
		fmt.Println(err)
	}
	var exist bool
	for result.Next() {
		result.Scan(&exist)
	}
	return exist
}
