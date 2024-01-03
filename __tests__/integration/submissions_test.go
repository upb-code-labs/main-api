package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSubmitSolutionToTestBlock(t *testing.T) {
	c := require.New(t)

	// ## Setup

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, status := CreateCourse("Submit solution to test block test - course")
	c.Equal(http.StatusCreated, status)

	// Create a laboratory
	laboratoryCreationResponse, status := CreateLaboratory(cookie, map[string]interface{}{
		"name":         "Submit solution to test block test - laboratory",
		"course_uuid":  courseUUID,
		"opening_date": "2023-12-01T08:00",
		"due_date":     "3023-12-01T00:00",
	})
	c.Equal(http.StatusCreated, status)
	laboratoryUUID := laboratoryCreationResponse["uuid"].(string)

	// Get the supported languages
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
		blockName:      "Update test block test - block",
		cookie:         cookie,
		testFile:       zipFile,
	})
	c.Equal(http.StatusCreated, status)
	testBlockUUID := blockCreationResponse["uuid"].(string)

	// Add a student to the course
	invitationCode, code := GetInvitationCode(courseUUID)
	c.Equal(http.StatusOK, code)
	c.NotEmpty(invitationCode)

	_, code = AddStudentToCourse(invitationCode)
	c.Equal(http.StatusOK, code)

	// ## Test
	zipFile, err = GetSampleSubmissionArchive()
	c.Nil(err)

	// Login as a student
	w, r = PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredStudentEmail,
		"password": registeredStudentPass,
	})
	router.ServeHTTP(w, r)
	cookie = w.Result().Cookies()[0]

	// Submit the solution
	submissionResponse, status := SubmitSolutionToTestBlock(&SubmitSToTestBlockUtilsDTO{
		blockUUID: testBlockUUID,
		cookie:    cookie,
		file:      zipFile,
	})

	c.Equal(http.StatusCreated, status)
	c.NotEmpty(submissionResponse["uuid"])
}
