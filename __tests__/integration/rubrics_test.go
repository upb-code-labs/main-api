package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateRubric(t *testing.T) {
	c := require.New(t)

	testCases := []GenericTestCase{
		{
			// Short username
			Payload: map[string]interface{}{
				"name": "a",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			// Valid data
			Payload: map[string]interface{}{
				"name": "Rubric 1",
			},
			ExpectedStatusCode: http.StatusCreated,
		},
	}

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Run test cases
	for _, testCase := range testCases {
		response, status := CreateRubric(cookie, testCase.Payload)
		c.Equal(testCase.ExpectedStatusCode, status)

		if testCase.ExpectedStatusCode == http.StatusCreated {
			c.NotEmpty(response["uuid"])
			c.Equal(testCase.Payload["name"], response["name"])
			c.NotEmpty(response["message"])
		}
	}
}

func CreateRubric(cookie *http.Cookie, payload map[string]interface{}) (response map[string]interface{}, status int) {
	w, r := PrepareRequest("POST", "/api/v1/rubrics", payload)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return ParseJsonResponse(w.Body), w.Code
}
