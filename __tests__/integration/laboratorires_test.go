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
				"due_date":     "2023-12-01T12:00",
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
