package classroom

import (
	connections "attsys/connections"
	keygen "attsys/keygen"
	"attsys/models"
	"encoding/base64"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"math"
	"net/http"

	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func convertToBase64String(ImagePath string) string {
	bytes, err := ioutil.ReadFile(ImagePath)
	if err != nil {
		fmt.Println("Error ! Image path is bad")
		return ""
	}
	var base64Encoding string
	mimeType := http.DetectContentType(bytes)
	switch mimeType {
	case "image/jpeg":
		base64Encoding += "data:image/jpeg;base64,"
	case "image/png":
		base64Encoding += "data:image/png;base64,"
	}

	// Append the base64 encoded output
	base64Encoding += base64.StdEncoding.EncodeToString(bytes)

	// Print the full base64 representation of the image
	return base64Encoding
}

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
	query := fmt.Sprintf(`select classroom_id,department_id,course_id,from_date,to_date from classrooms
	where teacher_id = '%s'`, TeacherId)
	result, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	var temp models.Classroom
	var listOfClassrooms []models.Classroom
	for result.Next() {
		result.Scan(&temp.ClassroomId, &temp.DepartmentId, &temp.CourseId, &temp.From, &temp.To)
		if temp.To.Sub(time.Now().UTC()) > 0 {
			listOfClassrooms = append(listOfClassrooms, temp)
		} else {
			fmt.Println(temp)
			fmt.Println("Expired classroom!")
		}
	}
	return listOfClassrooms
}

func GetSessionsOfClassroom(ClassroomId int) []models.Session {
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

func CreateUnqiueSession(newSession models.Session) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return false
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
		return false
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
			return false
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
			fmt.Println("Error inserting students into attendance!")
			fmt.Println(err)
			return false
		}

	} else {
		return false
	}
	return true

}

func IsAnySessionActive(ClassroomId int) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return false
	}
	defer db.Close()

	query := fmt.Sprintf(`select exists(
		select 1 from sessions 
		where classroom_id = %d and session_status = 'ACTIVE' or session_status = 'WAITING'
	)`, ClassroomId)

	result, err := db.Query(query)
	if err != nil {
		fmt.Println("Error occurred during checking if session already exists")
		fmt.Println(err)
		return true
	}
	var exist bool
	for result.Next() {
		result.Scan(&exist)
	}
	return exist
}

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

func GetSessionDetails(SessionId int) models.SessionDashBoardDetails {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
	}
	defer db.Close()
	query := fmt.Sprintf(`select d.department_name,c.course_name,concat(t.firstname,' ',t.lastname) as teachername,session_date,start_time,end_time,session_status,reviewed
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

	var Session models.SessionDashBoardDetails
	for result.Next() {
		fmt.Println(result.Scan(&Session.DepartmentName, &Session.CourseName,
			&Session.TeacherName, &Session.SessionDetails.Date, &Session.SessionDetails.Start_time,
			&Session.SessionDetails.End_time, &Session.SessionDetails.Status, &Session.SessionDetails.Reviewed))
	}

	query = fmt.Sprintf(`select session_key from keygen where session_id = %d`, SessionId)
	result, err = db.Query(query)
	if err != nil {
		fmt.Println("Error retrieving teacher session dash board details")
		return Session
	}
	for result.Next() {
		fmt.Println(result.Scan(&Session.SessionDetails.SessionKey))
	}

	fmt.Println(Session)
	return Session
}
func GetStudentOfSessionDetails(SessionId int) []models.AttendanceDetails {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
	}
	defer db.Close()

	var Students []models.AttendanceDetails
	query := fmt.Sprintf(`select concat(firstname,' ',lastname) as studentname, 
		picture,regnumber, attendance1,attendance1_fp ,attendance2,attendance2_fp, 
		attendance3,attendance3_fp, 
		ispresent from (
		select * from attendance
		left join attendance_image_table
		using(session_id,regnumber,classroom_id)
		where session_id = %d) as s1
		left join students
		using(regnumber);
	`, SessionId)

	result, err := db.Query(query)

	if err != nil {
		fmt.Println("Error retrieving student attending session details")
		return Students
	}

	for result.Next() {
		var temp models.AttendanceDetails
		var attfp1 string
		var attfp2 string
		var attfp3 string
		var fairImagePath string

		err := result.Scan(&temp.StudentName, &fairImagePath, &temp.Regnumber,
			&temp.Attendance1.Time, &attfp1,
			&temp.Attendance2.Time, &attfp2,
			&temp.Attendance3.Time, &attfp3,
			&temp.IsPresent)

		fmt.Println(err)
		fmt.Println("Student is present:", temp.IsPresent)

		temp.FairImage = template.URL(convertToBase64String(fairImagePath))

		temp.Attendance1.ImageFilePath = template.URL(convertToBase64String(attfp1))
		temp.Attendance2.ImageFilePath = template.URL(convertToBase64String(attfp2))
		temp.Attendance3.ImageFilePath = template.URL(convertToBase64String(attfp3))

		Students = append(Students, temp)

	}
	return Students
}

func IsSessionReviewed(SessionId int) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return true
	}
	defer db.Close()
	query := fmt.Sprintf("select reviewed from sessions where session_id = %d", SessionId)
	result, _ := db.Query(query)
	var isReviewed bool
	for result.Next() {
		result.Scan(&isReviewed)
	}
	return isReviewed
}

func ReviewAndSetAttendance(ClassroomId int, SessionId int, Attendance []models.ReviewAttendance) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return false
	}
	defer db.Close()
	JsonAttendance, err := json.Marshal(Attendance)
	if err != nil {
		fmt.Println("Error marshaling attedance")
		return false
	}

	query := fmt.Sprintf(`
		call set_attendance(%d,%d,'%s')
	`, ClassroomId, SessionId, string(JsonAttendance))
	fmt.Println(query)
	_, err = db.Query(query)
	if err == nil {
		query = fmt.Sprintf(`update sessions
				 set reviewed = true
				 where session_id = %d
		`, SessionId)
		_, err = db.Query(query)
		if err != nil {
			fmt.Println(err)
		}
		return err == nil
	} else {
		fmt.Println(err)
	}
	return false

}

func InsertStudentIntoClassroom(ClassroomId int, Regnumber string) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return false
	}
	defer db.Close()
	query := fmt.Sprintf(`insert into classroom_attendees
		values(%d,'%s')
	`, ClassroomId, Regnumber)
	_, err = db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	if IsAnySessionActive(ClassroomId) {
		query = fmt.Sprintf(`
			call insert_into_active_session(%d,'%s')
		`, ClassroomId, Regnumber)
		_, err = db.Query(query)
		fmt.Println(query)
		if err != nil {
			fmt.Println(err)
			RemoveStudentFromClassroom(ClassroomId, Regnumber)
		}
		return err == nil
	}
	return true
}
func RemoveStudentFromClassroom(ClassroomId int, Regnumber string) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return false
	}
	defer db.Close()
	query := fmt.Sprintf(`delete from classroom_attendees
	where classroom_id = %d and regnumber = '%s'`, ClassroomId, Regnumber)
	_, err = db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}

func GetStudentsOfClassroom(ClassroomId int) []models.StudentsDetails {

	var Students []models.StudentsDetails
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return Students
	}
	defer db.Close()
	query := fmt.Sprintf(`select concat(firstname,' ',lastname) as Studentname,regnumber,email,picture from (
		select * from classroom_attendees
		left join students
		using(regnumber)
		where classroom_id = %d
	) as s;
	`, ClassroomId)

	result, err := db.Query(query)
	if err == nil {
		for result.Next() {
			var temp models.StudentsDetails
			var ImagePath string
			err := result.Scan(&temp.Studentname, &temp.Regnumber, &temp.Email, &ImagePath)
			fmt.Println(err)
			fmt.Println(temp)
			temp.Image = template.URL(convertToBase64String(ImagePath))

			Students = append(Students, temp)
		}
	}
	return Students

}

func GetCourses() []models.Course {
	var Courses []models.Course
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return Courses
	}
	defer db.Close()
	query := `select course_id,course_name from courses`
	result, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error getting courses")
	}
	for result.Next() {
		var temp models.Course
		result.Scan(&temp.CourseId, &temp.CourseName)
		Courses = append(Courses, temp)
	}
	return Courses
}

func GetDepartments() []models.Department {
	var Departments []models.Department
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return Departments
	}
	defer db.Close()
	query := `select department_id,department_name from departments`
	result, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Error getting courses")
	}
	for result.Next() {
		var temp models.Department
		result.Scan(&temp.DepartmentId, &temp.DepartmentName)
		Departments = append(Departments, temp)
	}
	return Departments
}
