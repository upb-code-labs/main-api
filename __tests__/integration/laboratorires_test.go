package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateLaboratory(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, status := CreateCourse("Create laboratory test - course")
	c.Equal(http.StatusCreated, status)

	// Define tests cases
	testCases := []GenericTestCase{
		{
			Payload: map[string]interface{}{
				"name":         ".",
				"course_uuid":  "not a uuid",
				"opening_date": "not a date",
				"due_date":     "not a date",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Payload: map[string]interface{}{
				"name":         "Create laboratory test - laboratory",
				"course_uuid":  courseUUID,
				"opening_date": "2023-12-01T08:00",
				"due_date":     "3023-12-01T12:00",
			},
			ExpectedStatusCode: http.StatusCreated,
		},
	}

	// Run tests
	for _, tc := range testCases {
		_, status := CreateLaboratory(cookie, tc.Payload)
		c.Equal(tc.ExpectedStatusCode, status)
	}
}

func TestGetLaboratoryByUUID(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, status := CreateCourse("Get laboratory by uuid test - course")
	c.Equal(http.StatusCreated, status)

	// Create a laboratory
	laboratoryName := "Get laboratory by uuid test - laboratory"
	laboratoryOpeningDate := "2023-12-01T08:00"
	laboratoryDueDate := "2023-12-01T12:00"

	laboratoryCreationResponse, status := CreateLaboratory(cookie, map[string]interface{}{
		"name":         laboratoryName,
		"course_uuid":  courseUUID,
		"opening_date": laboratoryOpeningDate,
		"due_date":     laboratoryDueDate,
	})
	laboratoryUUID := laboratoryCreationResponse["uuid"].(string)
	c.Equal(http.StatusCreated, status)

	// Define tests cases
	testCases := []GenericTestCase{
		{
			Payload: map[string]interface{}{
				"laboratory_uuid": "4e2ba78e-a8f0-4312-b4a7-e8c6933029b8",
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Payload: map[string]interface{}{
				"laboratory_uuid": "not a uuid",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Payload: map[string]interface{}{
				"laboratory_uuid": laboratoryUUID,
			},
			ExpectedStatusCode: http.StatusOK,
		},
	}

	// Run tests
	for _, tc := range testCases {
		getLaboratoryResponse, status := GetLaboratoryByUUID(cookie, tc.Payload["laboratory_uuid"].(string))
		c.Equal(tc.ExpectedStatusCode, status)

		if tc.ExpectedStatusCode == http.StatusOK {
			// Validate string fields
			c.Equal(laboratoryUUID, getLaboratoryResponse["uuid"])
			c.Equal(laboratoryName, getLaboratoryResponse["name"])
			c.Nil(getLaboratoryResponse["rubric_uuid"])
			c.Contains(getLaboratoryResponse["opening_date"], laboratoryOpeningDate)
			c.Contains(getLaboratoryResponse["due_date"], laboratoryDueDate)

			// Validate blocks fields
			c.Equal(0, len(getLaboratoryResponse["markdown_blocks"].([]interface{})))
			c.Equal(0, len(getLaboratoryResponse["test_blocks"].([]interface{})))
		}
	}
}

func TestUpdateLaboratory(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, status := CreateCourse("Update laboratory test - course")
	c.Equal(http.StatusCreated, status)

	// Create a laboratory
	initialLaboratoryName := "Update laboratory test - laboratory"
	laboratoryOpeningDate := "2023-12-01T08:00"
	laboratoryDueDate := "2023-12-01T12:00"

	laboratoryCreationResponse, status := CreateLaboratory(cookie, map[string]interface{}{
		"name":         initialLaboratoryName,
		"course_uuid":  courseUUID,
		"opening_date": laboratoryOpeningDate,
		"due_date":     laboratoryDueDate,
	})
	laboratoryUUID := laboratoryCreationResponse["uuid"].(string)
	c.Equal(http.StatusCreated, status)

	// Create a rubric
	rubricName := "Update laboratory test - rubric"
	rubricCreationResponse, status := CreateRubric(cookie, map[string]interface{}{
		"name": rubricName,
	})
	rubricUUID := rubricCreationResponse["uuid"].(string)
	c.Equal(http.StatusCreated, status)

	// Define tests cases
	updatedLaboratoryName := "Update laboratory test - laboratory updated"
	testCases := []GenericTestCase{
		{
			Payload: map[string]interface{}{
				"laboratory_uuid": "ea21f0a2-713f-427a-94d4-f541281fd654",
				"rubric_uuid":     rubricUUID,
				"name":            updatedLaboratoryName,
				"opening_date":    laboratoryOpeningDate,
				"due_date":        laboratoryDueDate,
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Payload: map[string]interface{}{
				"laboratory_uuid": "not a uuid",
				"rubric_uuid":     rubricUUID,
				"name":            updatedLaboratoryName,
				"opening_date":    laboratoryOpeningDate,
				"due_date":        laboratoryDueDate,
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Payload: map[string]interface{}{
				"laboratory_uuid": laboratoryUUID,
				"rubric_uuid":     rubricUUID,
				"name":            updatedLaboratoryName,
				"opening_date":    laboratoryOpeningDate,
				"due_date":        laboratoryDueDate,
			},
			ExpectedStatusCode: http.StatusNoContent,
		},
	}

	// Run tests
	for _, tc := range testCases {
		_, status := UpdateLaboratory(cookie, tc.Payload["laboratory_uuid"].(string), tc.Payload)
		c.Equal(tc.ExpectedStatusCode, status)
	}

	// Validate laboratory update
	getLaboratoryResponse, status := GetLaboratoryByUUID(cookie, laboratoryUUID)
	c.Equal(http.StatusOK, status)
	c.Equal(updatedLaboratoryName, getLaboratoryResponse["name"])
	c.Equal(rubricUUID, getLaboratoryResponse["rubric_uuid"])
	c.Contains(getLaboratoryResponse["opening_date"], laboratoryOpeningDate)
	c.Contains(getLaboratoryResponse["due_date"], laboratoryDueDate)
	c.Equal(0, len(getLaboratoryResponse["markdown_blocks"].([]interface{})))
	c.Equal(0, len(getLaboratoryResponse["test_blocks"].([]interface{})))
}

func TestCreateMarkdownBlock(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, status := CreateCourse("Create markdown block test - course")
	c.Equal(http.StatusCreated, status)

	// Create a laboratory
	laboratoryCreationResponse, status := CreateLaboratory(cookie, map[string]interface{}{
		"name":         "Create markdown block test - laboratory",
		"course_uuid":  courseUUID,
		"opening_date": "2023-12-01T08:00",
		"due_date":     "3023-12-01T12:00",
	})
	laboratoryUUID := laboratoryCreationResponse["uuid"].(string)
	c.Equal(http.StatusCreated, status)

	// Define tests cases
	testCases := []GenericTestCase{
		{
			Payload: map[string]interface{}{
				"laboratory_uuid": "0e5890ce-bd0d-422b-bc1e-a1cfc7f152e9",
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Payload: map[string]interface{}{
				"laboratory_uuid": "not a uuid",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Payload: map[string]interface{}{
				"laboratory_uuid": laboratoryUUID,
			},
			ExpectedStatusCode: http.StatusCreated,
		},
	}

	// Run tests
	var createdBlockUUID string
	for _, tc := range testCases {
		response, status := CreateMarkdownBlock(cookie, tc.Payload["laboratory_uuid"].(string))
		c.Equal(tc.ExpectedStatusCode, status)

		if tc.ExpectedStatusCode == http.StatusCreated {
			createdBlockUUID = response["uuid"].(string)
		}
	}

	// Validate the block was created
	getLaboratoryResponse, status := GetLaboratoryByUUID(cookie, laboratoryUUID)
	c.Equal(http.StatusOK, status)
	c.Equal(1, len(getLaboratoryResponse["markdown_blocks"].([]interface{})))

	block := getLaboratoryResponse["markdown_blocks"].([]interface{})[0].(map[string]interface{})
	c.Equal(createdBlockUUID, block["uuid"])
	c.Equal("", block["content"])
	c.EqualValues(1, block["index"])
}

func TestCreateTestBlock(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, status := CreateCourse("Create test block test - course")
	c.Equal(http.StatusCreated, status)

	// Create a laboratory
	laboratoryCreationResponse, status := CreateLaboratory(cookie, map[string]interface{}{
		"name":         "Create test block test - laboratory",
		"course_uuid":  courseUUID,
		"opening_date": "2023-12-01T08:00",
		"due_date":     "3023-12-01T00:00",
	})
	laboratoryUUID := laboratoryCreationResponse["uuid"].(string)
	c.Equal(http.StatusCreated, status)

	// Get a supported language from the supported languages list
	language := GetFirstSupportedLanguage(cookie)
	languageUUID := language["uuid"].(string)

	// Open `.zip` file from the data folder
	zipFile, err := GetSampleTestsArchive()
	c.Nil(err)

	// Send the request
	response, _ := CreateTestBlock(&CreateTestBlockUtilsDTO{
		laboratoryUUID: laboratoryUUID,
		languageUUID:   languageUUID,
		blockName:      "Create test block test - block",
		cookie:         cookie,
		testFile:       zipFile,
	})

	// Validate the response
	// c.Equal(http.StatusCreated, status)
	c.Contains(response, "uuid")
}

func TestGetStudentsProgres(t *testing.T) {
	c := require.New(t)

	// ## Prepare

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, status := CreateCourse("Get students progress test - course")
	c.Equal(http.StatusCreated, status)

	// Create a laboratory
	laboratoryCreationResponse, status := CreateLaboratory(cookie, map[string]interface{}{
		"name":         "Get students progress test - laboratory",
		"course_uuid":  courseUUID,
		"opening_date": "2023-12-01T08:00",
		"due_date":     "3023-12-01T00:00",
	})
	laboratoryUUID := laboratoryCreationResponse["uuid"].(string)
	c.Equal(http.StatusCreated, status)

	// Add a student to the course
	invitationCode, code := GetInvitationCode(courseUUID)
	c.Equal(http.StatusOK, code)
	c.NotEmpty(invitationCode)

	_, code = AddStudentToCourse(invitationCode)
	c.Equal(http.StatusOK, code)

	// Get languages list
	languagesResponse, status := GetSupportedLanguages(cookie)
	c.Equal(http.StatusOK, status)

	languages := languagesResponse["languages"].([]interface{})
	c.Greater(len(languages), 0)

	firstLanguage := languages[0].(map[string]interface{})
	firstLanguageUUID := firstLanguage["uuid"].(string)

	// Create a test block
	zipFile, err := GetSampleTestsArchive()
	c.Nil(err)

	blockCreationResponse, status := CreateTestBlock(&CreateTestBlockUtilsDTO{
		laboratoryUUID: laboratoryUUID,
		languageUUID:   firstLanguageUUID,
		blockName:      "Get students progress test - block",
		cookie:         cookie,
		testFile:       zipFile,
	})
	c.Equal(http.StatusCreated, status)
	testBlockUUID := blockCreationResponse["uuid"].(string)

	// Login as a student
	w, r = PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredStudentEmail,
		"password": registeredStudentPass,
	})
	router.ServeHTTP(w, r)
	cookie = w.Result().Cookies()[0]

	// Submit a solution
	zipFile, err = GetSampleSubmissionArchive()
	c.Nil(err)

	_, status = SubmitSolutionToTestBlock(&SubmitSToTestBlockUtilsDTO{
		blockUUID: testBlockUUID,
		file:      zipFile,
		cookie:    cookie,
	})
	c.Equal(http.StatusCreated, status)

	// ## Test

	// Login as a teacher again
	w, r = PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie = w.Result().Cookies()[0]

	// Get the students progress
	response, status := GetStudentsProgressInLaboratory(laboratoryUUID, cookie)
	c.Equal(http.StatusOK, status)

	totalTestBlocks := response["total_test_blocks"].(float64)
	studentsProgress := response["students_progress"].([]interface{})
	c.Equal(1, int(totalTestBlocks))
	c.Equal(1, len(studentsProgress))

	studentProgress := studentsProgress[0].(map[string]interface{})
	studentPendingSubmissionsCount := studentProgress["pending_submissions"].(float64)
	studentRunningSubmissionsCount := studentProgress["running_submissions"].(float64)
	studentFailingSubmissionsCount := studentProgress["failing_submissions"].(float64)
	StudentSuccessSubmissionsCount := studentProgress["success_submissions"].(float64)

	c.GreaterOrEqual(int(studentPendingSubmissionsCount), 0)
	c.GreaterOrEqual(int(studentRunningSubmissionsCount), 0)
	c.GreaterOrEqual(int(studentFailingSubmissionsCount), 0)
	c.GreaterOrEqual(int(StudentSuccessSubmissionsCount), 0)
}
