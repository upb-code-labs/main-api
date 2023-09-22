package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterStudent(t *testing.T) {
	c := require.New(t)

	// Register route
	router.POST("/accounts/students", accountsControllers.HandleRegisterStudent)

	testCases := []TestCase{
		{
			Payload: map[string]interface{}{
				"full_name":        "John Doe",
				"email":            "Not an email",
				"institutional_id": "Not numeric",
				"password":         "short",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Payload: map[string]interface{}{
				"full_name":        "John Doe",
				"email":            "john.doe.2020@upb.edu.co",
				"institutional_id": "000486314",
				"password":         "john/password/2023",
			},
			ExpectedStatusCode: http.StatusCreated,
		},
		{
			// Same email
			Payload: map[string]interface{}{
				"full_name":        "John Doe",
				"email":            "john.doe.2020@upb.edu.co",
				"institutional_id": "000634814",
				"password":         "john/password/2023",
			},
			ExpectedStatusCode: http.StatusConflict,
		},
		{
			// Same institutional_id
			Payload: map[string]interface{}{
				"full_name":        "John Doe",
				"email":            "john.doe.2023@upb.edu.co",
				"institutional_id": "000486314",
				"password":         "john/password/2023",
			},
			ExpectedStatusCode: http.StatusConflict,
		},
	}

	for _, testCase := range testCases {
		w, r := PerformRequest("POST", "/accounts/students", testCase.Payload)
		router.ServeHTTP(w, r)
		c.Equal(testCase.ExpectedStatusCode, w.Code)
	}
}
