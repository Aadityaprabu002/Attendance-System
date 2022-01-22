package classroom

import (
	adminss "attsys/admin/backend/ss"
	student_classroom "attsys/classroom/backend/ss"
	connections "attsys/connections"
	key "attsys/env"
	"attsys/models"
	teacher "attsys/teacher/backend"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(key.GetSecretKey()))

func CreateOrAppendClassRoom(TeacherId string, Classroom models.Classroom) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic("failed to establish connection with sql")
	}
	defer db.Close()
	query := fmt.Sprintf(`insert into classrooms (teacher_id,department_id,course_id,from_date,to_date)
	values('%s','%s','%s','%s','%s');`, TeacherId, Classroom.DepartmentId, Classroom.CourseId, Classroom.From.Format("2006-01-02"), Classroom.To.Format("2006-01-02"))

	query = strings.TrimSpace(query)
	fmt.Println(query)
	_, err = db.Exec(query)
	if err != nil {
		fmt.Println("Error creating a classroom")
		fmt.Println(err)
		return false
	}
	return true
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	msg := models.Htmlresponse{
		Response: "",
		Status:   0,
	}
	if teacher.IsLogged(r) {
		session, _ := store.Get(r, "teacher")
		if r.Method == "GET" {
			fmt.Println("GET")
			tmp, _ := template.ParseFiles("classroom/frontend/ts/dashboard.html")
			data := models.ClassroomTableData{
				Classrooms: GetClassrooms(session.Values["TEACHER_ID"].(string)),
				Courses:    GetCourses(),
				Department: GetDepartments(),
			}
			fmt.Println("list of classrooms:")
			fmt.Print(data)
			tmp.Execute(w, data)

		} else if r.Method == "POST" {
			fmt.Println("POST")
			var params models.Classroom
			err := json.NewDecoder(r.Body).Decode(&params)
			if err != nil {
				fmt.Println(err)
			}
			if CreateOrAppendClassRoom(session.Values["TEACHER_ID"].(string), params) {
				msg.Status = 1
			} else {
				msg.Response = "Failed creating a classroom"
			}
			json.NewEncoder(w).Encode(msg)
		}
	} else {
		http.Redirect(w, r, "/teacher/signin", http.StatusSeeOther)
	}

}

func ClassroomDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	session, err := store.Get(r, "teacher")
	if err != nil {
		fmt.Println(err)
	}

	if teacher.IsLogged(r) {
		if r.Method == "GET" {
			fmt.Println("GET")
			tmp, _ := template.ParseFiles("classroom/frontend/ts/classroomdashboard.html")
			params := mux.Vars(r)
			ClassroomId, err := strconv.Atoi(params["ClassroomId"])
			if err != nil {
				fmt.Println("bad code for ClassroomId")
				http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
				return
			}
			TeacherId := session.Values["TEACHER_ID"].(string)

			if !isAuthenticClassroom(TeacherId, ClassroomId) {
				fmt.Println("Invalid ClassroomId for the logged teacher")
				http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
				return
			}
			Sessions := GetSessionsOfClassroom(ClassroomId)
			Students := GetStudentsOfClassroom(ClassroomId)

			ClassroomTableDetails := models.ClassroomTableDetails{
				Sessions: Sessions,
				Students: Students,
			}
			tmp.Execute(w, ClassroomTableDetails)

		} else if r.Method == "POST" {
			fmt.Println("POST")
			temp := mux.Vars(r)
			ClassroomId, err := strconv.Atoi(temp["ClassroomId"])
			if err != nil {
				fmt.Println("bad code for ClassroomId")
				return
			}

			msg := models.Htmlresponse{
				Status:   0,
				Response: "",
			}
			var params models.Session
			err = json.NewDecoder(r.Body).Decode(&params)
			if err != nil {
				fmt.Println("Error decoding post parameters")
				fmt.Println(err)
				return

			}
			params.ClassroomId = ClassroomId
			loc, _ := time.LoadLocation("Asia/Kolkata")
			params.Start_time = params.Start_time.In(loc)
			params.End_time = params.End_time.In(loc)

			fmt.Println("New session details:", params)

			if err != nil {
				fmt.Println("Error obtaining post values!")
				fmt.Println(err)
			}
			if IsAnySessionActive(ClassroomId) {
				msg.Response = "Already you have an active or waiting session"
			} else if params.Start_time.After(params.End_time) {
				msg.Response = "Start time greater than end time!"
			} else if params.End_time.Sub(params.Start_time) > time.Duration(4*time.Hour) {
				msg.Response = "Session time greater than 4 hours!"

			} else if params.End_time.Sub(params.Start_time) < time.Duration(20*time.Minute) {
				msg.Response = "Session time less than 20 minutes!"
			} else {
				status := CreateUnqiueSession(params)
				if status {
					msg.Status = 1
					msg.Response = "New Session Created!"

				} else {
					msg.Response = "Error creating session!"
				}
			}
			json.NewEncoder(w).Encode(msg)
		}
	} else {
		http.Redirect(w, r, "/teacher/signin", http.StatusSeeOther)
	}

}

func SessionDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if teacher.IsLogged(r) {
		if r.Method == "GET" {
			fmt.Println("GET")
			params := mux.Vars(r)
			ClassroomId, err := strconv.Atoi(params["ClassroomId"])
			if err != nil {
				fmt.Println("bad code for ClassroomId")
				http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
				return
			}
			SessionId, err := strconv.Atoi(params["SessionId"])
			if err != nil {
				fmt.Println("Bad code for SessionId")
				http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
				return
			}
			session, _ := store.Get(r, "teacher")
			TeacherId := session.Values["TEACHER_ID"].(string)

			if !isAuthenticSession(TeacherId, ClassroomId, SessionId) {
				fmt.Println("Invalid classroomId or sessionId for the logged teacher")
				http.Redirect(w, r, "/teacher/dashboard", http.StatusSeeOther)
			}

			var TeacherSessionDetails models.TeacherSessionDashBoardDetails
			TeacherSessionDetails.SessionDb = GetSessionDetails(SessionId)
			TeacherSessionDetails.Attendees = GetStudentOfSessionDetails(SessionId)
			tmp, _ := template.ParseFiles("classroom/frontend/ts/sessiondashboard.html")
			fmt.Println(tmp.Execute(w, TeacherSessionDetails))

		} else if r.Method == "POST" {

			fmt.Println("POST")
		}
	} else {
		http.Redirect(w, r, "/teacher/signin", http.StatusSeeOther)
	}
}

func PostAttendance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if teacher.IsLogged(r) {
		if r.Method == "POST" {
			res := models.Htmlresponse{
				Response: "Failed Posting attendance",
				Status:   0,
			}
			params := mux.Vars(r)
			ClassroomId, err := strconv.Atoi(params["ClassroomId"])
			if err != nil {
				fmt.Println("bad code for ClassroomId")
				res.Response = "bad code for ClassroomId"
				return
			}
			SessionId, err := strconv.Atoi(params["SessionId"])
			if err != nil {
				fmt.Println("Bad code for SessionId")
				res.Response = "bad code for SessionId"
				return
			}
			session, _ := store.Get(r, "teacher")
			TeacherId := session.Values["TEACHER_ID"].(string)
			if isAuthenticClassroom(TeacherId, ClassroomId) && isAuthenticSession(TeacherId, ClassroomId, SessionId) {
				if !IsSessionReviewed(SessionId) {
					Attendance := []models.ReviewAttendance{}
					err := json.NewDecoder(r.Body).Decode(&Attendance)
					if err != nil {
						fmt.Println("Error decoding student present list")
					} else {
						if ReviewAndSetAttendance(ClassroomId, SessionId, Attendance) {
							res.Response = "Attendance Reviewed and set successfully!"
							res.Status = 1
						} else {
							res.Response = "Failed Reviewing and Setting Attendance!"
						}

					}

				}

			}
			fmt.Println(res)
			json.NewEncoder(w).Encode(res)
		}
	}
}

func Handlestudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if teacher.IsLogged(r) {
		if r.Method == "POST" {
			res := models.Htmlresponse{
				Response: "Failed to add student!",
				Status:   0,
			}
			params := mux.Vars(r)
			ClassroomId, err := strconv.Atoi(params["ClassroomId"])
			if err != nil {
				fmt.Println("bad code for ClassroomId")
				res.Response = "bad code for ClassroomId"
				return
			}
			session, _ := store.Get(r, "teacher")
			TeacherId := session.Values["TEACHER_ID"].(string)
			if isAuthenticClassroom(TeacherId, ClassroomId) {
				Student := struct {
					Regnumber string `json:"regnumber"`
					Code      int    `json:"code"`
				}{}
				err := json.NewDecoder(r.Body).Decode(&Student)

				if err == nil {
					if adminss.IsStudentExist(Student.Regnumber) {
						switch Student.Code {
						case 1:
							if !student_classroom.IsStudentBelongsToClassroom(Student.Regnumber, ClassroomId) {
								if adminss.GetStudentAccountStatus(Student.Regnumber) == 2 {
									if InsertStudentIntoClassroom(ClassroomId, Student.Regnumber) {
										res.Response = "Successflly added student! " + Student.Regnumber
										res.Status = 1
									}
								} else {
									res.Response = "Student not completed registering"
								}
							} else {
								res.Response = "Student already present in the classroom"
							}
						case 2:
							if student_classroom.IsStudentBelongsToClassroom(Student.Regnumber, ClassroomId) {
								if RemoveStudentFromClassroom(ClassroomId, Student.Regnumber) {
									res.Response = "Successflly removed student!" + Student.Regnumber
									res.Status = 1
								} else {
									res.Response = "Student does not belong to the classroom or not registered"
								}
							} else {
								res.Response = "Student not present in the classroom"
							}
						}
					} else {
						res.Response = "Not a valid student"
					}

				} else {
					fmt.Println("Error decoding the request")
				}
			}

			json.NewEncoder(w).Encode(res)
		}
	}
}
