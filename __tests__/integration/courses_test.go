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
	RegisterStudentAccountPayload := requests.RegisterUserRequest{
		FullName:        "Jeffrey Richardson",
		Email:           "jeffrey.richardson.2020@upb.edu.co",
		InstitutionalId: "000345678",
		Password:        "jeffrey/password/2023",
	}
	code = RegisterStudentAccount(RegisterStudentAccountPayload)
	c.Equal(201, code)

	// Login with the student
	w, r = PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    RegisterStudentAccountPayload.Email,
		"password": RegisterStudentAccountPayload.Password,
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

type InvitationCodeTestCase struct {
	CourseUUID         string
	ExpectedStatusCode int
}

func TestGetCourseByUUID(t *testing.T) {
	c := require.New(t)

	// Create a course
	courseName := "Course [Test Get Course By UUID]"
	courseUUID, code := CreateCourse(courseName)
	c.Equal(http.StatusCreated, code)
	c.NotEmpty(courseUUID)

	testCases := []InvitationCodeTestCase{
		{
			CourseUUID:         "not-valid",
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			// Non-existent course
			CourseUUID:         "3febe413-d8cc-4d77-961a-cba1a4eaa64e",
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
		response, code := GetCourseByUUID(cookie, testCase.CourseUUID)
		c.Equal(testCase.ExpectedStatusCode, code)

		if code == http.StatusOK {
			c.Equal(courseName, response["name"])
			c.Equal(courseUUID, response["uuid"])
			c.NotEmpty(response["color"])
		}
	}
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

func TestGetCourses(t *testing.T) {
	c := require.New(t)

	// --- 1. Try with a student account ---
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredStudentEmail,
		"password": registeredStudentPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]
	response, code := GetCoursesUserIsEnrolledIn(cookie)

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
	response, code = GetCoursesUserIsEnrolledIn(cookie)

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
	response, code := GetCoursesUserIsEnrolledIn(cookie)
	c.Equal(http.StatusOK, code)
	c.Equal(1, len(response["hidden_courses"].([]interface{})))

	// Assertion> Show the course
	json, code = ToggleCourseVisibility(cookie, courseUUID)
	c.Equal(http.StatusOK, code)
	c.True(json["visible"].(bool))

	// Get the student courses
	response, code = GetCoursesUserIsEnrolledIn(cookie)
	c.Equal(http.StatusOK, code)
	c.Equal(0, len(response["hidden_courses"].([]interface{})))
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

func TestEnrollStudentToCourse(t *testing.T) {
	c := require.New(t)

	// Register an student
	RegisterStudentAccountPayload := requests.RegisterUserRequest{
		FullName:        "Karl Ivica",
		Email:           "karl.ivica.2020@upb.edu.co",
		InstitutionalId: "000814593",
		Password:        "karl/password/2023",
	}
	code := RegisterStudentAccount(RegisterStudentAccountPayload)
	c.Equal(http.StatusCreated, code)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Get the student UUID
	response, code := SearchStudentsByFullName(cookie, "Karl Ivica")
	students := response["students"].([]interface{})
	c.Equal(http.StatusOK, code)
	c.Equal(1, len(students))
	student_uuid := students[0].(map[string]interface{})["uuid"].(string)

	// Create a course
	courseUUID, code := CreateCourse("Course [Test Enroll Student]")
	c.Equal(http.StatusCreated, code)

	enrollTestCases := []GenericTestCase{
		{
			Payload: map[string]interface{}{
				// Non-valid student
				"student_uuid": "not-valid",
				"course_uuid":  courseUUID,
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Payload: map[string]interface{}{
				"student_uuid": student_uuid,
				// Non-valid course
				"course_uuid": "not-valid",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Payload: map[string]interface{}{
				"student_uuid": student_uuid,
				// Non-existent course
				"course_uuid": "7eb30f08-f097-4ff9-b760-2d692adda73a",
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Payload: map[string]interface{}{
				"student_uuid": student_uuid,
				"course_uuid":  courseUUID,
			},
			ExpectedStatusCode: http.StatusNoContent,
		},
		{
			Payload: map[string]interface{}{
				"student_uuid": student_uuid,
				"course_uuid":  courseUUID,
			},
			ExpectedStatusCode: http.StatusConflict,
		},
	}

	for _, testCase := range enrollTestCases {
		_, code := EnrollSTudentToCourse(cookie, testCase.Payload["course_uuid"].(string), testCase.Payload["student_uuid"].(string))
		c.Equal(testCase.ExpectedStatusCode, code)
	}

	// Get the enrolled students
	getEnrolledTestCases := []GenericTestCase{
		{
			Payload: map[string]interface{}{
				"course_uuid": "not-valid",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Payload: map[string]interface{}{
				// Non-existent course
				"course_uuid": "7eb30f08-f097-4ff9-b760-2d692adda73a",
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Payload: map[string]interface{}{
				"course_uuid": courseUUID,
			},
			ExpectedStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range getEnrolledTestCases {
		response, code := GetStudentsEnrolledInCourse(cookie, testCase.Payload["course_uuid"].(string))
		c.Equal(testCase.ExpectedStatusCode, code)

		if code == http.StatusOK {
			students := response["students"].([]interface{})
			c.Equal(1, len(students))

			// Assert the student fields
			student := students[0].(map[string]interface{})
			c.Equal(student_uuid, student["uuid"])
			c.Equal(RegisterStudentAccountPayload.FullName, student["full_name"])
			c.Equal(RegisterStudentAccountPayload.InstitutionalId, student["institutional_id"])
			c.Equal(true, student["is_active"])
		}
	}
}

func TestGetCourseLaboratories(t *testing.T) {
	c := require.New(t)

	// Create a course
	courseUUID, code := CreateCourse("Get course laboratories test - course")
	c.Equal(http.StatusCreated, code)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create two laboratories
	openLaboratoryName := "Get course laboratories test - open laboratory"
	openLaboratoryOpeningDate := "2023-11-01T08:00"
	openLaboratoryDueDate := "2023-11-07T00:00"
	openLaboratoryCreationResponse, code := CreateLaboratory(cookie, map[string]interface{}{
		"name":         openLaboratoryName,
		"course_uuid":  courseUUID,
		"opening_date": openLaboratoryOpeningDate,
		"due_date":     openLaboratoryDueDate,
	})
	c.Equal(http.StatusCreated, code)
	openLaboratoryUUID := openLaboratoryCreationResponse["uuid"].(string)

	futureLaboratoryName := "Get course laboratories test - future laboratory"
	_, code = CreateLaboratory(cookie, map[string]interface{}{
		"name":         futureLaboratoryName,
		"course_uuid":  courseUUID,
		"opening_date": "3023-11-01T08:00",
		"due_date":     "3023-11-07T00:00",
	})
	c.Equal(http.StatusCreated, code)

	// ## Teacher test cases
	teacherTestCases := []GenericTestCase{
		GenericTestCase{
			Payload: map[string]interface{}{
				courseUUID: "not a uuid",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				courseUUID: courseUUID,
			},
			ExpectedStatusCode: http.StatusOK,
		},
	}

	for _, tx := range teacherTestCases {
		response, code := GetCourseLaboratories(cookie, tx.Payload[courseUUID].(string))
		c.Equal(tx.ExpectedStatusCode, code)

		if code == http.StatusOK {
			laboratories := response["laboratories"].([]interface{})
			c.Equal(2, len(laboratories))

			// Assert the laboratories fields
			for _, laboratory := range laboratories {
				laboratory := laboratory.(map[string]interface{})
				c.NotEmpty(laboratory["uuid"])
				c.NotEmpty(laboratory["name"])
				c.NotEmpty(laboratory["opening_date"])
				c.NotEmpty(laboratory["due_date"])
			}
		}
	}

	// ## Student test cases
	// Enroll the student in the course
	invitationCode, code := GetInvitationCode(courseUUID)
	c.Equal(http.StatusOK, code)

	_, code = AddStudentToCourse(invitationCode)
	c.Equal(http.StatusOK, code)

	// Login as a student
	w, r = PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredStudentEmail,
		"password": registeredStudentPass,
	})
	router.ServeHTTP(w, r)
	cookie = w.Result().Cookies()[0]

	studentTestCases := []GenericTestCase{
		GenericTestCase{
			Payload: map[string]interface{}{
				courseUUID: courseUUID,
			},
			ExpectedStatusCode: http.StatusOK,
		},
	}

	for _, tx := range studentTestCases {
		response, code := GetCourseLaboratories(cookie, tx.Payload[courseUUID].(string))
		c.Equal(tx.ExpectedStatusCode, code)

		if code == http.StatusOK {
			laboratories := response["laboratories"].([]interface{})
			c.Equal(1, len(laboratories))

			// Assert the laboratories fields
			laboratory := laboratories[0].(map[string]interface{})
			c.Equal(openLaboratoryUUID, laboratory["uuid"])
			c.Equal(openLaboratoryName, laboratory["name"])
			c.Contains(laboratory["opening_date"], openLaboratoryOpeningDate)
			c.Contains(laboratory["due_date"], openLaboratoryDueDate)
		}
	}
}

func TestSetStudentStatus(t *testing.T) {
	c := require.New(t)

	// 1. Create a course
	courseUUID, code := CreateCourse("Set student status test - course")
	c.Equal(http.StatusCreated, code)

	// 2. Get the course invitation code
	invitationCode, code := GetInvitationCode(courseUUID)
	c.Equal(http.StatusOK, code)

	// 3. Add a student to the course
	_, code = AddStudentToCourse(invitationCode)
	c.Equal(http.StatusOK, code)

	// 4. Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// 5. Get the student uuid
	students, code := GetStudentsEnrolledInCourse(cookie, courseUUID)
	c.Equal(http.StatusOK, code)
	c.Equal(1, len(students["students"].([]interface{})))
	studentUUID := students["students"].([]interface{})[0].(map[string]interface{})["uuid"].(string)

	// 6. Set the student status
	code = SetStudentStatus(&SetStudentStatusUtilsDTO{
		CourseUUID:  courseUUID,
		StudentUUID: studentUUID,
		IsActive:    false,
		Cookie:      cookie,
	})
	c.Equal(http.StatusNoContent, code)

	// 7. Get the student status
	students, code = GetStudentsEnrolledInCourse(cookie, courseUUID)
	c.Equal(http.StatusOK, code)
	c.Equal(1, len(students["students"].([]interface{})))
	student := students["students"].([]interface{})[0].(map[string]interface{})
	c.Equal(false, student["is_active"])
}
