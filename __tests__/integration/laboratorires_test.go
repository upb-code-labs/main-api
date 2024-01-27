package integration

import (
	"net/http"
	"testing"

	"github.com/gabriel-vasile/mimetype"
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
				"opening_date": defaultLaboratoryOpeningDate,
				"due_date":     defaultLaboratoryDueDate,
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

	laboratoryCreationResponse, status := CreateLaboratory(cookie, map[string]interface{}{
		"name":         laboratoryName,
		"course_uuid":  courseUUID,
		"opening_date": defaultLaboratoryOpeningDate,
		"due_date":     defaultLaboratoryDueDate,
	})
	laboratoryUUID := laboratoryCreationResponse["uuid"].(string)
	c.Equal(http.StatusCreated, status)

	// Define tests cases
	randomUUID := "4e2ba78e-a8f0-4312-b4a7-e8c6933029b8"
	testCases := []GenericTestCase{
		{
			Payload: map[string]interface{}{
				"laboratory_uuid": randomUUID,
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
		getLaboratoryResponse, status := GetLaboratoryByUUID(
			cookie,
			tc.Payload["laboratory_uuid"].(string),
		)
		c.Equal(tc.ExpectedStatusCode, status)

		getLaboratoryInformationResponse, status := GetLaboratoryInformationByUUID(
			cookie,
			tc.Payload["laboratory_uuid"].(string),
		)
		c.Equal(tc.ExpectedStatusCode, status)

		if tc.ExpectedStatusCode == http.StatusOK {
			// ## Validate laboratory request
			// Validate string fields
			c.Equal(laboratoryUUID, getLaboratoryResponse["uuid"])
			c.Equal(laboratoryName, getLaboratoryResponse["name"])
			c.Nil(getLaboratoryResponse["rubric_uuid"])
			c.Equal(defaultLaboratoryOpeningDateUTC, getLaboratoryResponse["opening_date"])
			c.Equal(defaultLaboratoryDueDateUTC, getLaboratoryResponse["due_date"])

			// Validate blocks fields
			c.Equal(0, len(getLaboratoryResponse["markdown_blocks"].([]interface{})))
			c.Equal(0, len(getLaboratoryResponse["test_blocks"].([]interface{})))

			// ## Validate laboratory information request
			c.Equal(laboratoryUUID, getLaboratoryInformationResponse["uuid"])
			c.Equal(laboratoryName, getLaboratoryInformationResponse["name"])
			c.Nil(getLaboratoryInformationResponse["rubric_uuid"])
			c.Equal(defaultLaboratoryOpeningDateUTC, getLaboratoryInformationResponse["opening_date"])
			c.Equal(defaultLaboratoryDueDateUTC, getLaboratoryInformationResponse["due_date"])
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
	defaultLaboratoryOpeningDate := defaultLaboratoryOpeningDate
	defaultLaboratoryDueDate := defaultLaboratoryDueDate

	laboratoryCreationResponse, status := CreateLaboratory(cookie, map[string]interface{}{
		"name":         initialLaboratoryName,
		"course_uuid":  courseUUID,
		"opening_date": defaultLaboratoryOpeningDate,
		"due_date":     defaultLaboratoryDueDate,
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
				"opening_date":    defaultLaboratoryOpeningDate,
				"due_date":        defaultLaboratoryDueDate,
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		{
			Payload: map[string]interface{}{
				"laboratory_uuid": "not a uuid",
				"rubric_uuid":     rubricUUID,
				"name":            updatedLaboratoryName,
				"opening_date":    defaultLaboratoryOpeningDate,
				"due_date":        defaultLaboratoryDueDate,
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Payload: map[string]interface{}{
				"laboratory_uuid": laboratoryUUID,
				"rubric_uuid":     rubricUUID,
				"name":            updatedLaboratoryName,
				"opening_date":    defaultLaboratoryOpeningDate,
				"due_date":        defaultLaboratoryDueDate,
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
	c.Equal(defaultLaboratoryOpeningDateUTC, getLaboratoryResponse["opening_date"])
	c.Equal(defaultLaboratoryDueDateUTC, getLaboratoryResponse["due_date"])
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
		"opening_date": defaultLaboratoryOpeningDate,
		"due_date":     defaultLaboratoryDueDate,
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

	// ## Preparation
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
		"opening_date": defaultLaboratoryOpeningDate,
		"due_date":     defaultLaboratoryDueDate,
	})
	laboratoryUUID := laboratoryCreationResponse["uuid"].(string)
	c.Equal(http.StatusCreated, status)

	// Get a supported language from the supported languages list
	language := GetFirstSupportedLanguage(cookie)
	languageUUID := language["uuid"].(string)

	// Open `.zip` file from the data folder
	zipFile, err := GetSampleTestsArchive()
	c.Nil(err)

	// ## Tests: Create test block
	// Send the request
	response, _ := CreateTestBlock(&CreateTestBlockUtilsDTO{
		laboratoryUUID: laboratoryUUID,
		languageUUID:   languageUUID,
		blockName:      "Create test block test - block",
		cookie:         cookie,
		testFile:       zipFile,
	})
	createdTestBlockUUID := response["uuid"].(string)
	c.Equal(http.StatusCreated, status)
	c.NotEmpty(createdTestBlockUUID)

	// ## Tests: Download test archive
	testsArchiveBytes, status := GetTestsArchive(createdTestBlockUUID, cookie)
	c.Equal(http.StatusOK, status)

	// Check the length of the response
	c.Greater(len(testsArchiveBytes), 0)

	// Check the MIMETYPE
	mtype := mimetype.Detect(testsArchiveBytes)
	c.Equal("application/zip", mtype.String())
}

func TestGetStudentsProgress(t *testing.T) {
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
		"opening_date": defaultLaboratoryOpeningDate,
		"due_date":     defaultLaboratoryDueDate,
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

	testBlockName := "Get students progress test - block"
	blockCreationResponse, status := CreateTestBlock(&CreateTestBlockUtilsDTO{
		laboratoryUUID: laboratoryUUID,
		languageUUID:   firstLanguageUUID,
		blockName:      testBlockName,
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

	// ### Get the progress of all the students in the laboratory
	response, status := GetStudentsProgressInLaboratory(laboratoryUUID, cookie)
	c.Equal(http.StatusOK, status)

	totalTestBlocks := response["total_test_blocks"].(float64)
	studentsProgress := response["students_progress"].([]interface{})
	c.Equal(1, int(totalTestBlocks))
	c.Equal(1, len(studentsProgress))

	studentProgress := studentsProgress[0].(map[string]interface{})

	studentUUID := studentProgress["student_uuid"].(string)
	studentFullName := studentProgress["student_full_name"].(string)

	c.NotEmpty(studentUUID)
	c.NotEmpty(studentFullName)

	studentPendingSubmissionsCount := studentProgress["pending_submissions"].(float64)
	studentRunningSubmissionsCount := studentProgress["running_submissions"].(float64)
	studentFailingSubmissionsCount := studentProgress["failing_submissions"].(float64)
	StudentSuccessSubmissionsCount := studentProgress["success_submissions"].(float64)

	c.GreaterOrEqual(int(studentPendingSubmissionsCount), 0)
	c.GreaterOrEqual(int(studentRunningSubmissionsCount), 0)
	c.GreaterOrEqual(int(studentFailingSubmissionsCount), 0)
	c.GreaterOrEqual(int(StudentSuccessSubmissionsCount), 0)

	// ### Get the progress of a student in the laboratory
	response, status = GetProgressOfStudentInLaboratory(
		laboratoryUUID,
		studentUUID,
		cookie,
	)
	c.Equal(http.StatusOK, status)

	totalTestBlocks = response["total_test_blocks"].(float64)
	c.Equal(1, int(totalTestBlocks))

	submissions := response["submissions"].([]interface{})
	c.Equal(1, len(submissions))

	submission := submissions[0].(map[string]interface{})
	c.NotEmpty(submission["uuid"])
	c.NotEmpty(submission["archive_uuid"])
	c.Equal(testBlockName, submission["test_block_name"])
	c.Contains([]string{"pending", "running", "ready"}, submission["status"])
	c.Contains([]bool{true, false}, submission["is_passing"])
}
