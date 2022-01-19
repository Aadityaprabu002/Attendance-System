package classroom

import (
	admin "attsys/admin/backend/ss"
	key "attsys/env"
	models "attsys/models"
	student "attsys/student/backend"

	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(key.GetSecretKey()))

func Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if student.IsLogged(r) {
		session, _ := store.Get(r, "student")
		Regnumber := session.Values["REG_NUMBER"].(string)
		status := admin.GetStudentAccountStatus(Regnumber)
		switch status {
		case 0:
			http.Redirect(w, r, "/student/signin", http.StatusSeeOther)
		case 1:
			fmt.Println("Student should completely register first!")
			http.Redirect(w, r, "/student/signin/complete_registration", http.StatusSeeOther)
		case 2:

			if session.Values["SESSION_KEY"] != nil {
				http.Redirect(w, r, "/student/dashboard/session", http.StatusSeeOther)
			} else {
				if r.Method == "GET" {
					tmp, _ := template.ParseFiles("classroom/frontend/ss/dashboard.html")
					temp := admin.GetStudentDetails(Regnumber)
					Student := models.StudentsDetails{
						Studentname: temp.Firstname + temp.Lastname,
						Regnumber:   temp.Regnumber,
						Image:       template.URL(convertToBase64String(temp.Image)),
						Email:       temp.Email,
					}
					tmp.Execute(w, Student)
				} else if r.Method == "POST" {
					var newJoinee models.Joinee
					msg := models.Htmlresponse{
						Response: "",
						Status:   0,
					}
					err := json.NewDecoder(r.Body).Decode(&newJoinee)
					if err != nil {
						fmt.Println("Error obtaining the session-key from the joinee")
					}
					if SessionId := isValidSessionKey(newJoinee.SessionKey); SessionId != 0 {
						Regnumber := session.Values["REG_NUMBER"].(string)
						if IsStudentBelongsToSession(Regnumber, SessionId) {
							session.Values["SESSION_KEY"] = newJoinee.SessionKey
							session.Save(r, w)
							msg.Status = 1
							fmt.Println("Student belongs to the session!!")
						} else {
							msg.Status = 0
							msg.Response = "Your account is not registered for this classroom"
						}
					} else {
						msg.Status = 0
						msg.Response = "Not a valid session key"
					}
					json.NewEncoder(w).Encode(msg)
				}
			}
		}
	} else {
		http.Redirect(w, r, "/student/signin", http.StatusSeeOther)
	}
}

func SessionDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if student.IsLogged(r) {
		session, _ := store.Get(r, "student")
		Regnumber := session.Values["REG_NUMBER"].(string)
		status := admin.GetStudentAccountStatus(Regnumber)
		switch status {
		case 0:
			http.Redirect(w, r, "/student/signin", http.StatusSeeOther)
		case 1:
			http.Redirect(w, r, "/student/signin/completeregistration", http.StatusSeeOther)
		case 2:
			if session.Values["SESSION_KEY"] != nil {
				if r.Method == "GET" {

					SessionId := isValidSessionKey(session.Values["SESSION_KEY"].(string))

					if SessionId == 0 {
						http.Redirect(w, r, "/student/dashboard", http.StatusSeeOther)
					}

					StudentSessionDetails, _ := GetSessionDetails(SessionId)
					tmp, _ := template.ParseFiles("classroom/frontend/ss/sessiondashboard.html")

					tmp.Execute(w, StudentSessionDetails)

				}
			} else {
				fmt.Println("SESSION KEY NOT SET!")
				http.Redirect(w, r, "/student/dashboard", http.StatusSeeOther)
			}
		}

	} else {
		http.Redirect(w, r, "/student/signin", http.StatusSeeOther)
	}

}
func SessionDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if student.IsLogged(r) {
		if r.Method == "GET" {
			session, _ := store.Get(r, "student")
			Regnumber := session.Values["REG_NUMBER"].(string)
			status := admin.GetStudentAccountStatus(Regnumber)
			if status == 2 {
				if session.Values["SESSION_KEY"] != nil {
					SessionId := isValidSessionKey(session.Values["SESSION_KEY"].(string))
					if SessionId != 0 {
						Date, StartTime, EndTime := GetSessionTimings(SessionId)

						StartTime = GetIndianTimeStamp(Date, StartTime)
						EndTime = GetIndianTimeStamp(Date, EndTime)

						PopUp1, PopUp2, PopUp3 := GetPopUpTimings(SessionId)

						PopUp1 = GetIndianTimeStamp(Date, PopUp1)
						PopUp2 = GetIndianTimeStamp(Date, PopUp2)
						PopUp3 = GetIndianTimeStamp(Date, PopUp3)

						res := models.StudentSessionTimerDetails{
							StartTime: StartTime,
							EndTime:   EndTime,
							PopUp1:    PopUp1,
							PopUp2:    PopUp2,
							PopUp3:    PopUp3,
						}
						json.NewEncoder(w).Encode(res)
					}
				}
			}
		}
	}
}

func PostAttendance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		res := models.Htmlresponse{
			Status:   0,
			Response: "Failed to post attendance!",
		}
		if student.IsLogged(r) {
			session, _ := store.Get(r, "student")
			Regnumber := session.Values["REG_NUMBER"].(string)
			status := admin.GetStudentAccountStatus(Regnumber)
			if status == 2 {
				SessionId := isValidSessionKey(session.Values["SESSION_KEY"].(string))
				if SessionId != 0 {
					Regnumber := session.Values["REG_NUMBER"].(string)
					if status := GetSessionStatus(SessionId); status == "ACTIVE" {
						params := models.PostAttendance{}
						err := json.NewDecoder(r.Body).Decode(&params)
						if err == nil {
							if InsertAttendance(Regnumber, SessionId, params) {
								res.Status = 1
								res.Response = "Attendance Posted Successfull!"
							}
						} else {
							fmt.Println("Error Decoding attendance request")
							fmt.Println(err)
						}
					}
				} else {
					res.Response = "Invalid Session to post attendance!"
				}
			}
		}
		json.NewEncoder(w).Encode(res)
	}
}
func EndSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	res := models.Htmlresponse{
		Response: "Failed Exiting Session",
		Status:   0,
	}
	if r.Method == "GET" {
		if student.IsLogged(r) {
			session, _ := store.Get(r, "student")
			Regnumber := session.Values["REG_NUMBER"].(string)
			status := admin.GetStudentAccountStatus(Regnumber)
			if status == 2 {
				if session.Values["SESSION_KEY"] != nil {
					session.Values["SESSION_KEY"] = nil
					session.Save(r, w)
					res.Status = 1
					res.Response = "Session Exited successfully!"
				}
			}
		}
		json.NewEncoder(w).Encode(res)
	}
}
