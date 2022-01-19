package classroom

import (
	connections "attsys/connections"
	"attsys/models"
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
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

func IsStudentBelongsToClassroom(Regnumber string, ClassroomId int) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return false
	}
	defer db.Close()
	query := fmt.Sprintf(`select exists (
		select 1 from classroom_attendees 
		where classroom_id = %d and regnumber = '%s'
	)`, ClassroomId, Regnumber)

	result, err := db.Query(query)
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

/*
Grabs the Regsiter number, SessionId and then checks
whether the student belongs to the classroom associated with the Session given
Returns true if student belongs otherwise false
*/
func IsStudentBelongsToSession(Regnumber string, SessionId int) bool {
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
	exists := IsStudentBelongsToClassroom(Regnumber, ClassroomId)
	return exists
}

func GetSessionDetails(SessionId int) (models.SessionDashBoardDetails, bool) {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	var data models.SessionDashBoardDetails
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
		result.Scan(&data.TeacherName, &data.SessionDetails.Date,
			&data.SessionDetails.Start_time, &data.SessionDetails.End_time,
			&data.SessionDetails.Status, &data.DepartmentName, &data.CourseName)
	}
	return data, true
}

func GetSessionStatus(SessionId int) string {
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

func GetSessionTimings(SessionId int) (time.Time, time.Time, time.Time) {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	var StartTime time.Time
	var EndTime time.Time
	var Date time.Time
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return Date, StartTime, EndTime
	}
	defer db.Close()

	query := fmt.Sprintf(`select session_date,start_time,end_time from sessions where session_id = %d`, SessionId)
	result, _ := db.Query(query)
	for result.Next() {
		result.Scan(&Date, &StartTime, &EndTime)
	}
	return Date, StartTime, EndTime
}

func GetPopUpTimings(SessionId int) (time.Time, time.Time, time.Time) {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	var PopUp1 time.Time
	var PopUp2 time.Time
	var PopUp3 time.Time
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return PopUp1, PopUp2, PopUp3
	}
	defer db.Close()
	query := fmt.Sprintf("select popup1,popup2,popup3 from keygen where session_id = %d", SessionId)
	result, _ := db.Query(query)
	for result.Next() {
		result.Scan(&PopUp1, &PopUp2, &PopUp3)
	}
	return PopUp1, PopUp2, PopUp3
}

func saveSessionStudentImageData(Regnumber string, SessionId int, ImageData string, ImageDataNumber int) string {
	fpath := fmt.Sprintf("../database/sessions/%d/%s/", SessionId, Regnumber)
	_ = os.MkdirAll(fpath, 0777)
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
	ImagePath := fmt.Sprintf("%sattendance%d.png", fpath, ImageDataNumber)
	ImageFile, err := os.OpenFile(ImagePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("Cannot open file")
		return ""
	}
	png.Encode(ImageFile, Image)
	return ImagePath

}

func InsertAttendance(Regnumber string, SessionId int, att models.PostAttendance) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println("failed to establish connection with sql")
		return false
	}
	defer db.Close()

	loc, _ := time.LoadLocation("Asia/Kolkata")
	att.Time = att.Time.In(loc) // Indian time

	attcolumns := [...]string{"attendance1", "attendance2", "attendance3"}
	query := fmt.Sprintf(`
		update attendance
		set %s = '%s'
		where regnumber = '%s' and session_id = %d
	`, attcolumns[att.AttNum-1], att.Time.Format("15:04:05"), Regnumber, SessionId)

	_, err = db.Query(query)
	if err != nil {
		fmt.Println("Error accepting attendance!")
		fmt.Println(err)
		return false
	}

	ImagePath := saveSessionStudentImageData(Regnumber, SessionId, att.Image, att.AttNum)
	attcolumns = [...]string{"attendance1_fp", "attendance2_fp", "attendance3_fp"}
	query = fmt.Sprintf(`
		update attendance_image_table
		set %s = '%s'
		where regnumber = '%s' and session_id = %d
	`, attcolumns[att.AttNum-1], ImagePath, Regnumber, SessionId)
	_, err = db.Query(query)
	if err != nil {
		fmt.Println("Error accepting attendance!")
		fmt.Println(err)
		return false
	}
	return true
}

func GetIndianTimeStamp(Date time.Time, Time time.Time) time.Time {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	INDtime := time.Date(
		Date.Year(),
		Date.Month(),
		Date.Day(),
		Time.Hour(),
		Time.Minute(),
		Time.Second(),
		0,
		loc,
	)
	return INDtime
}
