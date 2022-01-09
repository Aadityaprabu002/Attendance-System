package classroom

import (
	connections "attsys/connections"
	keygen "attsys/keygen"
	"attsys/models"
	"database/sql"
	"fmt"
	"strings"
	"time"
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

func getSessions(ClassroomId int) []models.Session {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
	}
	defer db.Close()
	query := fmt.Sprintf("select session_id,session_date,start_time,end_time,session_status from sessions where classroom_id = %d;", ClassroomId)
	result, err := db.Query(query)
	if err != nil {
		fmt.Println("Error retrieving session details for the classroom")
		fmt.Println(err)
	}
	var temp models.Session
	var listOfSessions []models.Session
	for result.Next() {
		result.Scan(&temp.SessionId, &temp.Date, &temp.Start_time, &temp.End_time, &temp.Status)
		listOfSessions = append(listOfSessions, temp)
	}
	return listOfSessions
}

func createUnqiueSession(newSession models.Session) (bool, int) {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return false, 0
	}
	defer db.Close()
	justNow := time.Now()
	sessionUniqueCode := fmt.Sprintf("%s%d%d%d", keygen.String(5), justNow.Hour(), justNow.Minute(), justNow.Second())
	fmt.Println(sessionUniqueCode)
	query := fmt.Sprintf("insert into sessions(session_date,start_time,end_time,classroom_id) values ('%s','%s','%s',%d);", newSession.End_time.Format("2006-01-02"), newSession.Start_time.Format("15:04:05"), newSession.End_time.Format("15:04:05"), newSession.ClassroomId)
	fmt.Println(query)
	_, err = db.Query(query)
	if err != nil {
		fmt.Println("Error creating a session")
		fmt.Println(err)
		return false, 0
	}
	query = fmt.Sprintf("select session_id from sessions where session_date = '%s' and start_time = '%s' and end_time = '%s';", newSession.End_time.Format("2006-01-02"), newSession.Start_time.Format("15:04:05"), newSession.End_time.Format("15:04:05"))
	fmt.Println(query)
	result, err := db.Query(query)
	if err != nil {
		fmt.Println("Error retrieving a session id")
		fmt.Println(err)
		return false, 0
	}
	var sid int
	for result.Next() {
		result.Scan(&sid)
	}
	if sid != 0 {
		query = fmt.Sprintf("insert into keygen(session_key,session_id) values('%s',%d)", sessionUniqueCode, sid)
		_, err = db.Query(query)
		if err != nil {
			fmt.Println("Error creating a session key")
			fmt.Println(err)
			return false, 0
		}
	}
	return true, sid

}

func checkForSession(ClassroomId int) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
	}
	defer db.Close()

	query := `select max(start_time) <= now()::time and now()::time < max(end_time) from sessions
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
