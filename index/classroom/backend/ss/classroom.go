package classroom

import (
	key "attsys/env"
	models "attsys/models"
	student "attsys/student/backend"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
	"time"

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
	if session.Values["SESSION_ID"] != nil {
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
					session.Values["SESSION_ID"] = SessionId
					EndTime := GetSessionEndTime(SessionId).String()
					fmt.Println(EndTime)
					session.Values["SESSION_END_TIME"] = EndTime
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

	if session.Values["SESSION_ID"] != nil {
		if r.Method == "GET" {
			SessionId := session.Values["SESSION_ID"].(int)
			StudentSessionDetails, _ := GetSessionDetails(SessionId)
			tmp, _ := template.ParseFiles("classroom/frontend/ss/sessiondashboard.html")
			tmp.Execute(w, StudentSessionDetails)
		} else if r.Method == "POST" {
			SessionId := session.Values["SESSION_ID"].(int)
			msg := models.Htmlresponse{
				Response: "",
				Status:   0,
			}
			status := IsSessionActive(SessionId)
			fmt.Println("Session Status:", status)
			switch status {
			case "CLOSED":
				session.Values["SESSION_ID"] = nil
				session.Values["SESSION_END_TIME"] = nil
				session.Save(r, w)
				msg.Status = 2
				msg.Response = "Session has been closed!"

			case "ACTIVE":
				fmt.Println(session.Values["SESSION_END_TIME"].(string))
				layout := "2006-01-02 15:04:05 -0700 MST"
				EndTime, _ := time.Parse(layout, session.Values["SESSION_END_TIME"].(string))
				CurrentTimeStr := "0000-01-01 " + time.Now().Format("15:04:05") + " +0000 UTC"
				CurrentTime, _ := time.Parse(layout, CurrentTimeStr)

				fmt.Println(EndTime, CurrentTime)
				msg.Status = 1
				msg.Response = GetSessionTimer(CurrentTime, EndTime)
				fmt.Println(msg.Response)

			case "WAITING":
				msg.Status = 0
				msg.Response = "Session not yet opened!"

			default:
				session.Values["SESSION_ID"] = nil
				session.Values["SESSION_END_TIME"] = nil
				session.Save(r, w)
				msg.Status = -1
				msg.Response = "Error!!! Failed to retrieve session status!"
			}
			json.NewEncoder(w).Encode(msg)
		}
	} else {
		fmt.Println("SESSION ID NOT SET!")
		http.Redirect(w, r, "/student/dashboard", http.StatusSeeOther)
	}

}
