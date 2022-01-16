package classroom

import (
	connections "attsys/connections"
	keygen "attsys/keygen"
	"attsys/models"
	"math"

	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func randRange(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max+1-min)
}
func GetPopUpTimings(StartTime time.Time, EndTime time.Time) (time.Time, time.Time, time.Time) {
	diff := EndTime.Sub(StartTime)
	Minutes := int(diff.Minutes())
	Interval := Minutes / 3
	var MultiplicationFactor = randRange(3, int(math.Max(3, float64(Interval-3))))
	popAt := time.Duration(int(time.Minute) * MultiplicationFactor)
	firstPopUp := StartTime.Add(popAt)
	secondPopUp := StartTime.Add(time.Duration(int(time.Minute)*Interval) + popAt)
	thirdPopUp := StartTime.Add(time.Duration(int(time.Minute)*Interval*2) + popAt)
	return firstPopUp, secondPopUp, thirdPopUp

}
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
	query := fmt.Sprintf("insert into sessions(session_date,start_time,end_time,classroom_id) values ('%s','%s','%s',%d) returning session_id", newSession.Date.Format("2006-01-02"), newSession.Start_time.Format("15:04:05"), newSession.End_time.Format("15:04:05"), newSession.ClassroomId)
	fmt.Println(query)
	result, err := db.Query(query)
	if err != nil {
		fmt.Println("Error creating a session")
		fmt.Println(err)
		return false, 0
	}
	var sid int
	for result.Next() {
		result.Scan(&sid)
	}
	if sid != 0 {
		popup1, popup2, popup3 := GetPopUpTimings(newSession.Start_time, newSession.End_time)
		query = fmt.Sprintf("insert into keygen(session_key,session_id,popup1,popup2,popup3) values('%s',%d,'%s','%s','%s')", sessionUniqueCode, sid, popup1.Format("15:04:05"), popup2.Format("15:04:05"), popup3.Format("15:04:05"))
		_, err = db.Query(query)
		if err != nil {
			fmt.Println("Error creating a session key")
			fmt.Println(err)
			return false, 0
		}

		query = fmt.Sprintf(`
		insert into attendance
		select classroom_id, %d as session_id,regnumber from classroom_attendees
		where classroom_id = (
			select classroom_id from sessions 
			where session_id = %d
		)`, sid, sid)
		_, err = db.Query(query)
		if err != nil {
			fmt.Println("Error creating a session key")
			fmt.Println(err)
			return false, 0
		}

	} else {
		return false, 0
	}
	return true, sid

}

// func checkForSession(ClassroomId int) bool {
// 	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
// 	db, err := sql.Open("postgres", conn)
// 	if err != nil {
// 		fmt.Println("failed to establish connection with sql")
// 		return false
// 	}
// 	defer db.Close()

// 	query := `select max(start_time) <= now()::time and now()::time < max(end_time) from sessions
// where session_date = (
// select max(session_date) from sessions) and classroom_id = %d`

// 	query = fmt.Sprintf(query, ClassroomId)
// 	query = strings.TrimSpace(query)

// 	result, err := db.Query(query)
// 	if err != nil {
// 		fmt.Println("Error occurred during checking if session already exists")
// 		fmt.Println(err)
// 		return false
// 	}
// 	var exist bool
// 	for result.Next() {
// 		result.Scan(&exist)
// 	}
// 	return exist
// }

func isAuthenticClassroom(TeacherId string, ClassroomId int) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return false
	}
	defer db.Close()

	query := fmt.Sprintf(`select exists(
		select 1 from classrooms where classroom_id = '%d' and teacher_id = '%s'
		)`, ClassroomId, TeacherId)

	result, err := db.Query(query)
	if err != nil {
		fmt.Println("failed checking for authentic classroom")
		return false
	}
	var valid bool
	for result.Next() {
		result.Scan(&valid)
	}
	return valid
}

func isAuthenticSession(TeacherId string, ClassroomId int, SessionId int) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return false
	}
	defer db.Close()
	query := fmt.Sprintf(`select exists(
		select 1 from (
				select teacher_id,classroom_id,s.session_id from classrooms 
				inner join
				(
					select session_id,classroom_id from sessions
				)as s
				using(classroom_id)
			) as t
			where t.classroom_id = %d and t.session_id = %d and t.teacher_id = '%s'
		);
	`, ClassroomId, SessionId, TeacherId)
	result, err := db.Query(query)
	if err != nil {
		fmt.Println("failed checking for authentic session")
		return false
	}
	var valid bool
	for result.Next() {
		result.Scan(&valid)
	}
	return valid
}

func GetSessionDetails(SessionId int) models.TeacherSessionDashBoard {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
	}
	defer db.Close()
	query := fmt.Sprintf(`select d.department_name,c.course_name,concat(t.firstname,' ',t.lastname) as teachername,session_date,start_time,end_time,session_status 
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
	using(teacher_id)
	`, SessionId)
	result, err := db.Query(query)
	if err != nil {
		fmt.Println("Error retrieving teacher session dash board details")
	}
	var data models.TeacherSessionDashBoard
	for result.Next() {
		var temp models.Session
		result.Scan(&data.DepartmentName, &data.CourseName, &data.TeacherName, &temp.Date, &temp.Start_time, &temp.End_time, &data.SessionDetails.Status)
		data.SessionDetails.Date = temp.Date.Format("2006-01-02")
		data.SessionDetails.Start_time = temp.Start_time.Format("15:04:05")
		data.SessionDetails.End_time = temp.End_time.Format("15:04:05")
	}
	query = fmt.Sprintf(`select session_key from keygen where session_id = %d`, SessionId)
	result, err = db.Query(query)
	if err != nil {
		fmt.Println("Error retrieving teacher session dash board details")
		return data
	}
	for result.Next() {
		result.Scan(&data.SessionDetails.SessionKey)
	}

	query = fmt.Sprintf(`select concat(firstname,' ',lastname) as studentname,regnumber, attendance1,attendance1_fp ,attendance2,attendance2_fp, attendance3,attendance3_fp  from (
		select * from attendance
		left join attendance_image_table
		using(session_id,regnumber,classroom_id)
		where session_id = %d) as s1
		left join students
		using(regnumber);
	`, SessionId)
	result, err = db.Query(query)
	if err != nil {
		fmt.Println("Error retrieving teacher session dash board details")
		return data
	}
	var attendeesData models.AttendanceDetails
	var attendeesList []models.AttendanceDetails
	for result.Next() {
		var attTime1 time.Time
		var attTime2 time.Time
		var attTime3 time.Time
		result.Scan(&attendeesData.StudentName, &attendeesData.Regnumber,
			&attTime1, &attendeesData.Attendance1.ImageFilePath,
			&attTime2, &attendeesData.Attendance2.ImageFilePath,
			&attTime3, &attendeesData.Attendance3.ImageFilePath)

		attendeesData.Attendance1.PrettyTime = attTime1.Format("15:04:05")
		attendeesData.Attendance2.PrettyTime = attTime2.Format("15:04:05")
		attendeesData.Attendance3.PrettyTime = attTime3.Format("15:04:05")
		attendeesList = append(attendeesList, attendeesData)
	}
	data.Attendees = attendeesList
	return data
}
