package classroom

import (
	key "attsys/env"
	models "attsys/models"
	student "attsys/student/backend"

	"encoding/json"
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(key.GetSecretKey()))

func Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if !student.IsLogged(r) {
		http.Redirect(w, r, "/student/signin", http.StatusSeeOther)
	}
	session, err := store.Get(r, "student")
	if err != nil {
		fmt.Println("Error retrieving student session")
	}
	if session.Values["SESSION_KEY"] != nil {
		http.Redirect(w, r, "/student/dashboard/session", http.StatusSeeOther)
	} else {
		if r.Method == "GET" {
			fmt.Println("GET")
			tmp, _ := template.ParseFiles("classroom/frontend/ss/dashboard.html")
			tmp.Execute(w, nil)
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
				if IsStudentBelongsToClassroom(Regnumber, SessionId) {
					session.Values["SESSION_KEY"] = newJoinee.SessionKey
					fmt.Println(session.Save(r, w))
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

func SessionDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if !student.IsLogged(r) {
		http.Redirect(w, r, "/student/signin", http.StatusSeeOther)
	}

	session, err := store.Get(r, "student")
	if err != nil {
		fmt.Println("Error retrieving student session")
	}

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
func SessionDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if student.IsLogged(r) {
		if r.Method == "GET" {
			session, _ := store.Get(r, "student")
			SessionId := isValidSessionKey(session.Values["SESSION_KEY"].(string))
			if SessionId == 0 {
				return
			}
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

			// fmt.Println("Session And Pop up timing Details:", res)
			json.NewEncoder(w).Encode(res)
		}
	}
}

func PostAttendance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		if student.IsLogged(r) {
			session, _ := store.Get(r, "student")
			SessionId := isValidSessionKey(session.Values["SESSION_KEY"].(string))
			res := models.Htmlresponse{
				Status:   0,
				Response: "",
			}
			if SessionId == 0 {
				res.Response = "Failed to post attendance"
				json.NewEncoder(w).Encode(res)
				return
			}
			Regnumber := session.Values["REG_NUMBER"].(string)
			if status := GetSessionStatus(SessionId); status == "ACTIVE" {
				params := models.PostAttendance{}
				err := json.NewDecoder(r.Body).Decode(&params)
				if err != nil {
					fmt.Println("Error Decoding attendance request")
					fmt.Println(err)
					res.Response = "Failed to post attendance"
					json.NewEncoder(w).Encode(res)
					return
				}

				if InsertAttendance(Regnumber, SessionId, params) {
					res.Status = 1
					res.Response = "Attendance Posted Successfull!"
				} else {
					res.Response = "Failed to post attendance!"
				}
				json.NewEncoder(w).Encode(res)
			}
		}
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
			if session.Values["SESSION_KEY"] != nil {
				session.Values["SESSION_KEY"] = nil
				session.Save(r, w)
				res.Status = 1
				res.Response = "Session Exited successfully!"
			}
		}
		json.NewEncoder(w).Encode(res)
	}
}
