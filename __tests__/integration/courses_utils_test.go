package integration

import (
	"fmt"
	"net/http"
)

func CreateCourse(name string) (courseUUID string, statusCode int) {
	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create the course
	w, r = PrepareRequest("POST", "/api/v1/courses", map[string]interface{}{
		"name": name,
	})
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse["uuid"].(string), w.Code
}

func GetCourseByUUID(cookie *http.Cookie, courseUUID string) (response map[string]interface{}, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/courses/%s", courseUUID)
	w, r := PrepareRequest("GET", endpoint, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func GetInvitationCode(courseUUID string) (invitationCode string, statusCode int) {
	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Get the invitation code
	endpoint := fmt.Sprintf("/api/v1/courses/%s/invitation-code", courseUUID)
	w, r = PrepareRequest("GET", endpoint, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse["code"].(string), w.Code
}

func AddStudentToCourse(invitationCode string) (response map[string]interface{}, statusCode int) {
	// Login as a student
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredStudentEmail,
		"password": registeredStudentPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Join the course
	endpoint := fmt.Sprintf("/api/v1/courses/join/%s", invitationCode)
	w, r = PrepareRequest("POST", endpoint, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return ParseJsonResponse(w.Body), w.Code
}

func ToggleCourseVisibility(cookie *http.Cookie, courseUUID string) (response map[string]interface{}, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/courses/%s/visibility", courseUUID)
	w, r := PrepareRequest("PATCH", endpoint, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func GetCoursesUserIsEnrolledIn(cookie *http.Cookie) (response map[string]interface{}, statusCode int) {
	w, r := PrepareRequest("GET", "/api/v1/courses", nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func RenameCourse(cookie *http.Cookie, courseUUID string, name string) (statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/courses/%s/name", courseUUID)
	w, r := PrepareRequest("PATCH", endpoint, map[string]interface{}{
		"name": name,
	})
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return w.Code
}

func EnrollSTudentToCourse(cookie *http.Cookie, courseUUID string, studentUUID string) (response map[string]interface{}, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/courses/%s/students", courseUUID)
	w, r := PrepareRequest("POST", endpoint, map[string]interface{}{
		"student_uuid": studentUUID,
	})
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func GetStudentsEnrolledInCourse(cookie *http.Cookie, courseUUID string) (response map[string]interface{}, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/courses/%s/students", courseUUID)
	w, r := PrepareRequest("GET", endpoint, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func GetCourseLaboratories(cookie *http.Cookie, courseUUID string) (response map[string]interface{}, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/courses/%s/laboratories", courseUUID)
	w, r := PrepareRequest("GET", endpoint, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

type SetStudentStatusUtilsDTO struct {
	CourseUUID  string
	StudentUUID string
	IsActive    bool
	Cookie      *http.Cookie
}

func SetStudentStatus(dto *SetStudentStatusUtilsDTO) (statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/courses/%s/students/%s/status", dto.CourseUUID, dto.StudentUUID)
	w, r := PrepareRequest("PATCH", endpoint, map[string]interface{}{
		"is_active": dto.IsActive,
	})
	r.AddCookie(dto.Cookie)
	router.ServeHTTP(w, r)

	return w.Code
}
