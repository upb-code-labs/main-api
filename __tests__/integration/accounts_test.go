package integration

import (
	"net/http"
	"testing"

	"github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure/requests"
	"github.com/stretchr/testify/require"
)

type RegisterTestCase struct {
	Payload            map[string]interface{}
	ExpectedStatusCode int
}

func TestRegisterStudent(t *testing.T) {
	c := require.New(t)

	// Register route
	router.POST("/accounts/students", accountsControllers.HandleRegisterStudent)

	testCases := []RegisterTestCase{
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
		w, r := PrepareRequest("POST", "/accounts/students", testCase.Payload)
		router.ServeHTTP(w, r)
		c.Equal(testCase.ExpectedStatusCode, w.Code)
	}
}

func RegisterStudent(req requests.RegisterUserRequest) int {
	w, r := PrepareRequest("POST", "/accounts/students", map[string]interface{}{
		"full_name":        req.FullName,
		"email":            req.Email,
		"institutional_id": req.InstitutionalId,
		"password":         req.Password,
	})

	router.ServeHTTP(w, r)
	return w.Code
}

func TestRegisterAdmin(t *testing.T) {
	c := require.New(t)

	// Register route
	router.POST("/accounts/admins", accountsControllers.HandleRegisterAdmin)

	testCases := []RegisterTestCase{
		{
			Payload: map[string]interface{}{
				"full_name": "John Doe",
				"email":     "Not an email",
				"password":  "short",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Payload: map[string]interface{}{
				"full_name": "John Doe",
				"email":     "john.doe@gmail.com",
				"password":  "john/password/2023",
			},
			ExpectedStatusCode: http.StatusCreated,
		},
		{
			Payload: map[string]interface{}{
				"full_name": "John Doe",
				"email":     "john.doe@gmail.com",
				"password":  "john/password/2023",
			},
			ExpectedStatusCode: http.StatusConflict,
		},
	}

	for _, testCase := range testCases {
		w, r := PrepareRequest("POST", "/accounts/admins", testCase.Payload)
		router.ServeHTTP(w, r)
		c.Equal(testCase.ExpectedStatusCode, w.Code)
	}
}

func RegisterAdmin(req requests.RegisterAdminRequest) int {
	w, r := PrepareRequest("POST", "/accounts/admins", map[string]interface{}{
		"full_name": req.FullName,
		"email":     req.Email,
		"password":  req.Password,
	})

	router.ServeHTTP(w, r)
	return w.Code
}
