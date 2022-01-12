package classroom

import (
	connections "attsys/connections"
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

func isValidSessionKey(SessionKey string) int {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return 0
	}
	defer db.Close()
	query := fmt.Sprintf(`select session_id from keygen where session_key = '%s'`, SessionKey)
	result, err := db.Query(query)
	if err != nil {
		fmt.Println("Error while retrieving session id based on session key")
		return 0
	}
	var SessionId int
	for result.Next() {
		result.Scan(&SessionId)
	}
	return SessionId
}

/*
Grabs the Regsiter number, SessionId and then checks
whether the student belongs to the classroom associated with the Session given
Returns true if student belongs otherwise false
*/
func IsStudentBelongsToClassroom(Regnumber string, SessionId int) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(("failed to establish connection with sql"))
		return false
	}
	defer db.Close()
	query := fmt.Sprintf(`select classroom_id from sessions where session_id = %d`, SessionId)
	result, err := db.Query(query)
	if err != nil {
		fmt.Println("Error retrieving ClassroomId from sessions")
		return false
	}
	var ClassroomId int
	for result.Next() {
		result.Scan(&ClassroomId)
	}
	query = fmt.Sprintf(`select exists (
		select 1 from classroom_attendees 
		where classroom_id = %d and regnumber = '%s'
	)`, ClassroomId, Regnumber)
	result, err = db.Query(query)
	if err != nil {
		fmt.Println("Failed checking whether student exists in the classroom")
		fmt.Println(err)
		return false
	}
	var exists bool
	for result.Next() {
		result.Scan(&exists)
	}
	return exists
}

func GetSessionDetails(SessionId int) (models.StudentSessionDashBoard, bool) {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	var data models.StudentSessionDashBoard
	if err != nil {
		fmt.Println(("failed to establish connection with sql"))
		return data, false
	}
	defer db.Close()
	query := fmt.Sprintf(`select concat(t.firstname,' ',t.lastname) as teachername,session_date,start_time,end_time,session_status,d.department_name,c.course_name
	from (
			select * from classrooms
			inner join sessions
			using(classroom_id)
			where session_id = %d
	) as s1
	
	left join departments as d
	using (department_id)
	left join courses as c
	using (course_id)
	left join teachers as t
	using(teacher_id)`, SessionId)

	result, err := db.Query(query)
	if err != nil {
		fmt.Println("Cant retrieve Student Session details!")
		return data, false
	}

	for result.Next() {
		var temp models.Session
		result.Scan(&data.TeacherName, &temp.Date, &temp.Start_time, &temp.End_time, &data.SessionDetails.Status, &data.DepartmentName, &data.CourseName)
		data.SessionDetails.Date = temp.Date.Format("2006-01-02")
		data.SessionDetails.Start_time = temp.Start_time.Format("15:04:05")
		data.SessionDetails.End_time = temp.End_time.Format("15:04:05")
	}
	return data, true
}

func IsSessionActive(SessionId int) string {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return ""
	}
	defer db.Close()
	query := fmt.Sprintf(`select session_status from sessions where session_id = %d`, SessionId)
	result, err := db.Query(query)
	if err != nil {
		fmt.Println("Failed to retrieve session status")
		fmt.Println(err)
		return ""
	}
	var status string
	for result.Next() {
		result.Scan(&status)
	}
	return status
}

func GetSessionEndTime(SessionId int) time.Time {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	var EndTime time.Time
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return EndTime
	}
	defer db.Close()
	query := fmt.Sprintf(`select end_time from sessions where session_id = %d`, SessionId)
	result, _ := db.Query(query)
	for result.Next() {
		result.Scan(&EndTime)
	}
	return EndTime
}
func GetSessionTimer(CurrentTime time.Time, EndTime time.Time) string {
	diff := EndTime.Sub(CurrentTime)
	Hours := int(diff.Seconds()) / 3600
	Minutes := int(diff.Seconds()) / 60
	Seconds := int(diff.Seconds()) % 60
	return fmt.Sprintf("%d-%d-%d", Hours, Minutes, Seconds)
}
