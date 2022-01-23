package admin

import (
	classroom_teacher "attsys/classroom/backend/ts"
	"attsys/models"
	"encoding/json"
	"fmt"
	"net/http"
	"text/template"
)

func DepartmentAndCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Method == "GET" {
		tmp, _ := template.ParseFiles("admin/frontend/department_and_course/handle.html")
		res := struct {
			Courses    []models.Course
			Department []models.Department
		}{}
		res.Courses = classroom_teacher.GetCourses()
		res.Department = classroom_teacher.GetDepartments()
		tmp.Execute(w, res)
	}
}

func HandleDepartment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		params := struct {
			DepartmentId   string `json:"deptid"`
			DepartmentName string `json:"deptname"`
			Type           string `json:"action"`
		}{}
		json.NewDecoder(r.Body).Decode(&params)
		fmt.Println(params)
		res := models.Htmlresponse{
			Status:   0,
			Response: "Failed to add department",
		}

		if params.Type == "add" {
			if !IsDepartmentExist(params.DepartmentId) {
				var newDept models.Department
				newDept.DepartmentId = params.DepartmentId
				newDept.DepartmentName = params.DepartmentName
				if InsertDepartment(newDept) {
					res.Status = 1
					res.Response = "Successfully added department"
				}
			} else {

				res.Response = "Department already exists"
			}
		} else if params.Type == "delete" {
			if IsDepartmentExist(params.DepartmentId) {
				if DeleteDepartment(params.DepartmentId) {
					res.Status = 1
					res.Response = "Successfully department deleted"
				}
			} else {
				res.Response = "Department does not exists"
			}
		}
		json.NewEncoder(w).Encode(res)
	}

}

func HandleCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == "POST" {
		params := struct {
			CourseId   string `json:"courseid"`
			CourseName string `json:"coursename"`
			Type       string `json:"action"`
		}{}
		json.NewDecoder(r.Body).Decode(&params)
		fmt.Println(params)
		res := models.Htmlresponse{
			Status:   0,
			Response: "Failed to add Course",
		}

		if params.Type == "add" {
			if !IsCourseExist(params.CourseId) {
				var newCourse models.Course
				newCourse.CourseId = params.CourseId
				newCourse.CourseName = params.CourseName
				if InsertCourse(newCourse) {
					res.Status = 1
					res.Response = "Successfully added course"
				}
			} else {
				res.Response = "course already exist"
			}
		} else if params.Type == "delete" {
			if IsCourseExist(params.CourseId) {
				if DeleteCourse(params.CourseId) {
					res.Status = 1
					res.Response = "Successfully Course deleted"
				}
			} else {
				res.Response = "Successfully does not exist"
			}
		}
		json.NewEncoder(w).Encode(res)
	}

}
