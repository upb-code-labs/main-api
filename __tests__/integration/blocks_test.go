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
		"due_date":     "2023-12-01T12:00",
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
