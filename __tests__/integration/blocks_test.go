package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateMarkdownBlockContent(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, status := CreateCourse("Update markdown block content test - course")
	c.Equal(http.StatusCreated, status)

	// Create a laboratory
	laboratoryCreationResponse, status := CreateLaboratory(cookie, map[string]interface{}{
		"name":         "Update markdown block content test - laboratory",
		"course_uuid":  courseUUID,
		"opening_date": "2023-12-01T08:00",
		"due_date":     "3023-12-01T12:00",
	})
	c.Equal(http.StatusCreated, status)
	laboratoryUUID := laboratoryCreationResponse["uuid"].(string)

	// Create a markdown block
	blockCreationResponse, status := CreateMarkdownBlock(cookie, laboratoryUUID)
	c.Equal(http.StatusCreated, status)
	markdownBlockUUID := blockCreationResponse["uuid"].(string)

	// Define tests cases
	testCases := []GenericTestCase{
		{
			Payload: map[string]interface{}{
				"content": "",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Payload: map[string]interface{}{
				"content": "# Updated main title",
			},
			ExpectedStatusCode: http.StatusNoContent,
		},
	}

	// Run tests
	for _, tc := range testCases {
		_, status := UpdateMarkdownBlockContent(cookie, markdownBlockUUID, tc.Payload)
		c.Equal(tc.ExpectedStatusCode, status)
	}

	// Verify that the content was updated
	response, status := GetLaboratoryByUUID(cookie, laboratoryUUID)
	c.Equal(http.StatusOK, status)
	c.Equal(1, len(response["markdown_blocks"].([]interface{})))

	block := response["markdown_blocks"].([]interface{})[0].(map[string]interface{})
	c.Equal("# Updated main title", block["content"].(string))
}

func TestDeleteMarkdownBlock(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, status := CreateCourse("Delete markdown block test - course")
	c.Equal(http.StatusCreated, status)

	// Create a laboratory
	laboratoryCreationResponse, status := CreateLaboratory(cookie, map[string]interface{}{
		"name":         "Delete markdown block test - laboratory",
		"course_uuid":  courseUUID,
		"opening_date": "2023-12-01T08:00",
		"due_date":     "3023-12-01T12:00",
	})
	c.Equal(http.StatusCreated, status)
	laboratoryUUID := laboratoryCreationResponse["uuid"].(string)

	// Create a markdown block
	blockCreationResponse, status := CreateMarkdownBlock(cookie, laboratoryUUID)
	c.Equal(http.StatusCreated, status)
	markdownBlockUUID := blockCreationResponse["uuid"].(string)

	// Get the laboratory
	response, status := GetLaboratoryByUUID(cookie, laboratoryUUID)
	c.Equal(http.StatusOK, status)
	c.Equal(1, len(response["markdown_blocks"].([]interface{})))

	// Delete the markdown block
	_, status = DeleteMarkdownBlock(cookie, markdownBlockUUID)
	c.Equal(http.StatusNoContent, status)

	// Verify that the block was deleted
	response, status = GetLaboratoryByUUID(cookie, laboratoryUUID)
	c.Equal(http.StatusOK, status)
	c.Equal(0, len(response["markdown_blocks"].([]interface{})))
}

func TestUpdateTestBlock(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, status := CreateCourse("Update test block test - course")
	c.Equal(http.StatusCreated, status)

	// Create a laboratory
	laboratoryCreationResponse, status := CreateLaboratory(cookie, map[string]interface{}{
		"name":         "Update test block test - laboratory",
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

	// Update the test block
	newName := "Update test block test - block - updated"
	zipFile, err = GetSampleTestsArchive()
	c.Nil(err)

	_, status = UpdateTestBlock(&UpdateTestBlockUtilsDTO{
		blockUUID:    testBlockUUID,
		languageUUID: firstLanguageUUID,
		blockName:    newName,
		cookie:       cookie,
		testFile:     zipFile,
	})
	c.Equal(http.StatusNoContent, status)

	// Check that the test block data was updated
	laboratoryResponse, status := GetLaboratoryByUUID(cookie, laboratoryUUID)
	c.Equal(http.StatusOK, status)

	blocks := laboratoryResponse["test_blocks"].([]interface{})
	c.Equal(1, len(blocks))

	block := blocks[0].(map[string]interface{})
	c.Equal(newName, block["name"].(string))
	c.Equal(firstLanguageUUID, block["language_uuid"].(string))
}
