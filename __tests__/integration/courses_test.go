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
