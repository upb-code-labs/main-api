package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure/requests"
	"github.com/stretchr/testify/require"
)

func TestCreateCourse(t *testing.T) {
	c := require.New(t)

	testCases := []GenericTestCase{
		{
			Payload: map[string]interface{}{
				"name": "Course 1",
			},
			ExpectedStatusCode: http.StatusCreated,
		},
		{
			Payload: map[string]interface{}{
				"name": "a", // Short name
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
	}

	// --- 1. Try with a teacher account ---
	// Register a teacher
	registerTeacherPayload := requests.RegisterTeacherRequest{
		FullName: "Alayna Hartman",
		Email:    "alayna.hartman.2020@upb.edu.co",
		Password: "alayna/password/2023",
	}
	code := RegisterTeacherAccount(registerTeacherPayload)
	c.Equal(201, code)

	// Login with the teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registerTeacherPayload.Email,
		"password": registerTeacherPayload.Password,
	})
	router.ServeHTTP(w, r)
	hasCookie := len(w.Result().Cookies()) == 1
	c.True(hasCookie)
	cookie := w.Result().Cookies()[0]

	for _, testCase := range testCases {
		w, r = PrepareRequest("POST", "/api/v1/courses", testCase.Payload)
		r.AddCookie(cookie)
		router.ServeHTTP(w, r)

		jsonResponse := ParseJsonResponse(w.Body)
		c.Equal(testCase.ExpectedStatusCode, w.Code)

		// Check fields if the course was created
		if w.Code == http.StatusCreated {
			c.Equal(testCase.Payload["name"], jsonResponse["name"])
			c.NotEmpty(jsonResponse["uuid"])
			c.NotEmpty(jsonResponse["color"])
		}
	}

	// --- 2. Try with a non-teacher account ---
	// Register an student
	registerStudentPayload := requests.RegisterUserRequest{
		FullName:        "Jeffrey Richardson",
		Email:           "jeffrey.richardson.2020@upb.edu.co",
		InstitutionalId: "000345678",
		Password:        "jeffrey/password/2023",
	}
	code = RegisterStudent(registerStudentPayload)
	c.Equal(201, code)

	// Login with the student
	w, r = PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registerStudentPayload.Email,
		"password": registerStudentPayload.Password,
	})
	router.ServeHTTP(w, r)
	hasCookie = len(w.Result().Cookies()) == 1
	c.True(hasCookie)
	cookie = w.Result().Cookies()[0]

	for _, testCase := range testCases {
		w, r = PrepareRequest("POST", "/api/v1/courses", testCase.Payload)
		r.AddCookie(cookie)
		router.ServeHTTP(w, r)
		c.Equal(http.StatusForbidden, w.Code)
	}
}

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

type InvitationCodeTestCase struct {
	CourseUUID         string
	ExpectedStatusCode int
}

func TestGetInvitationCode(t *testing.T) {
	c := require.New(t)

	// Create a course
	courseUUID, code := CreateCourse("Course [Test Get Invitation Code]")
	c.Equal(http.StatusCreated, code)
	c.NotEmpty(courseUUID)

	testCases := []InvitationCodeTestCase{
		{
			CourseUUID:         "not-valid",
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			// Non-existent course
			CourseUUID:         "d1c80308-9e5e-42c9-a18b-7c2d2f78525e",
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			CourseUUID:         courseUUID,
			ExpectedStatusCode: http.StatusOK,
		},
	}

	// --- 1. Try with a teacher account ---
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	for _, testCase := range testCases {
		endpoint := fmt.Sprintf("/api/v1/courses/%s/invitation-code", testCase.CourseUUID)
		w, r = PrepareRequest("GET", endpoint, nil)
		r.AddCookie(cookie)
		router.ServeHTTP(w, r)

		c.Equal(testCase.ExpectedStatusCode, w.Code)
		if w.Code == http.StatusOK {
			jsonResponse := ParseJsonResponse(w.Body)
			c.NotEmpty(jsonResponse["code"])
		}
	}

	// --- 2. Try with a non-teacher account ---
	w, r = PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredStudentEmail,
		"password": registeredStudentPass,
	})
	router.ServeHTTP(w, r)
	cookie = w.Result().Cookies()[0]

	for _, testCase := range testCases {
		endpoint := fmt.Sprintf("/api/v1/courses/%s/invitation-code", testCase.CourseUUID)
		w, r = PrepareRequest("GET", endpoint, nil)
		r.AddCookie(cookie)
		router.ServeHTTP(w, r)
		c.Equal(http.StatusForbidden, w.Code)
	}
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

type JoinCourseTestCase struct {
	InvitationCode     string
	ExpectedStatusCode int
}

func TestJoinCourse(t *testing.T) {
	c := require.New(t)

	// Create a course
	courseUUID, code := CreateCourse("Course [Test Join Course]")
	c.Equal(http.StatusCreated, code)
	c.NotEmpty(courseUUID)

	// Get the invitation code
	invitationCode, code := GetInvitationCode(courseUUID)
	c.Equal(http.StatusOK, code)
	c.NotEmpty(invitationCode)

	testCases := []JoinCourseTestCase{
		{
			InvitationCode:     "a", // Invalid code
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			InvitationCode:     "abcdefghi", // Non-existent code
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			InvitationCode:     invitationCode,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			InvitationCode:     invitationCode, // Already joined
			ExpectedStatusCode: http.StatusConflict,
		},
	}

	// --- 1. Try with a student account ---

	for _, testCase := range testCases {
		response, code := AddStudentToCourse(testCase.InvitationCode)
		c.Equal(testCase.ExpectedStatusCode, code)

		// Check the response fields
		if code == http.StatusOK {
			c.NotEmpty(response["course"])
			c.NotEmpty(response["course"].(map[string]interface{})["uuid"])
			c.NotEmpty(response["course"].(map[string]interface{})["name"])
			c.NotEmpty(response["course"].(map[string]interface{})["color"])
		}
	}

	// --- 2. Try with a teacher account ---
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	for _, testCase := range testCases {
		endpoint := fmt.Sprintf("/api/v1/courses/join/%s", testCase.InvitationCode)
		w, r = PrepareRequest("POST", endpoint, nil)
		r.AddCookie(cookie)
		router.ServeHTTP(w, r)
		c.Equal(http.StatusForbidden, w.Code)
	}
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

func TestGetCourses(t *testing.T) {
	c := require.New(t)

	// --- 1. Try with a student account ---
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredStudentEmail,
		"password": registeredStudentPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]
	response, code := GetUserCourses(cookie)

	// Assertions
	c.Equal(http.StatusOK, code)
	assertGetCoursesResponse(c, response)

	// --- 2. Try with a teacher account ---
	w, r = PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie = w.Result().Cookies()[0]
	response, code = GetUserCourses(cookie)

	// Assertions
	c.Equal(http.StatusOK, code)
	assertGetCoursesResponse(c, response)
}

func assertGetCoursesResponse(c *require.Assertions, response map[string]interface{}) {
	c.NotEmpty(response["courses"])
	c.Empty(response["hidden_courses"])

	// Assert course fields
	courses := response["courses"].([]interface{})
	for _, course := range courses {
		course := course.(map[string]interface{})
		c.NotEmpty(course["uuid"])
		c.NotEmpty(course["name"])
		c.NotEmpty(course["color"])
	}
}

func GetUserCourses(cookie *http.Cookie) (response map[string]interface{}, statusCode int) {
	w, r := PrepareRequest("GET", "/api/v1/courses", nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func TestToggleCourseVisibility(t *testing.T) {
	c := require.New(t)

	// --- 1. Try with a student account ---
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredStudentEmail,
		"password": registeredStudentPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Assertion> Try with an invalid course
	_, code := ToggleCourseVisibility(cookie, "not-valid")
	c.Equal(http.StatusBadRequest, code)

	// Create a course
	courseUUID, code := CreateCourse("Course [Test Toggle Visibility]")
	c.Equal(http.StatusCreated, code)
	c.NotEmpty(courseUUID)

	// Assertion> Try to hide the course without being enrolled
	_, code = ToggleCourseVisibility(cookie, courseUUID)
	c.Equal(http.StatusForbidden, code)

	// Add a student to the course
	invitationCode, code := GetInvitationCode(courseUUID)
	c.Equal(http.StatusOK, code)
	c.NotEmpty(invitationCode)
	_, code = AddStudentToCourse(invitationCode)
	c.Equal(http.StatusOK, code)

	// Assertion> Hide the course
	json, code := ToggleCourseVisibility(cookie, courseUUID)
	c.Equal(http.StatusOK, code)
	c.False(json["visible"].(bool))

	// Get the student courses
	response, code := GetUserCourses(cookie)
	c.Equal(http.StatusOK, code)
	c.Equal(1, len(response["hidden_courses"].([]interface{})))

	// Assertion> Show the course
	json, code = ToggleCourseVisibility(cookie, courseUUID)
	c.Equal(http.StatusOK, code)
	c.True(json["visible"].(bool))

	// Get the student courses
	response, code = GetUserCourses(cookie)
	c.Equal(http.StatusOK, code)
	c.Equal(0, len(response["hidden_courses"].([]interface{})))
}

func ToggleCourseVisibility(cookie *http.Cookie, courseUUID string) (response map[string]interface{}, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/courses/%s/visibility", courseUUID)
	w, r := PrepareRequest("PATCH", endpoint, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func TestRenameCourse(t *testing.T) {
	c := require.New(t)

	testCases := []GenericTestCase{
		{
			Payload: map[string]interface{}{
				// Short name
				"name": "a",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Payload: map[string]interface{}{
				"name": "Competitive Programming",
			},
			ExpectedStatusCode: http.StatusNoContent,
		},
		{
			Payload: map[string]interface{}{
				// Same name
				"name": "Competitive Programming",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
	}

	// --- 1. Try with a teacher account ---
	// Create a course
	courseUUID, code := CreateCourse("Course [Test Rename Course]")
	c.Equal(http.StatusCreated, code)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	for _, testCase := range testCases {
		endpoint := fmt.Sprintf("/api/v1/courses/%s/name", courseUUID)
		w, r = PrepareRequest("PATCH", endpoint, testCase.Payload)
		r.AddCookie(cookie)
		router.ServeHTTP(w, r)
		c.Equal(testCase.ExpectedStatusCode, w.Code)
	}

	// Try with a non-valid course
	code = RenameCourse(cookie, "not-valid", "New Name")
	c.Equal(http.StatusBadRequest, code)

	// Try with a non-existent course
	code = RenameCourse(cookie, "ab41d891-8374-4eec-adef-5a129986b059", "New Name")
	c.Equal(http.StatusNotFound, code)

	// --- 2. Try with a teacher that does not own the course ---
	// Register a teacher
	registerTeacherPayload := requests.RegisterTeacherRequest{
		FullName: "Santeri Rasim",
		Email:    "santeri.rasim.2020@upb.edu.co",
		Password: "santeri/password/2023",
	}
	code = RegisterTeacherAccount(registerTeacherPayload)
	c.Equal(201, code)

	// Login with the teacher
	w, r = PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registerTeacherPayload.Email,
		"password": registerTeacherPayload.Password,
	})
	router.ServeHTTP(w, r)
	cookie = w.Result().Cookies()[0]

	code = RenameCourse(cookie, courseUUID, "New Name")
	c.Equal(http.StatusForbidden, code)
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
