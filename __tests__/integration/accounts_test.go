package integration

import (
	"net/http"
	"testing"

	"github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure/requests"
	"github.com/stretchr/testify/require"
)

func TestRegisterStudent(t *testing.T) {
	c := require.New(t)

	testCases := []GenericTestCase{
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
		w, r := PrepareRequest("POST", "/api/v1/accounts/students", testCase.Payload)
		router.ServeHTTP(w, r)
		c.Equal(testCase.ExpectedStatusCode, w.Code)
	}
}

func RegisterStudent(req requests.RegisterUserRequest) int {
	w, r := PrepareRequest("POST", "/api/v1/accounts/students", map[string]interface{}{
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

	// Login as an admin
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredAdminEmail,
		"password": registeredAdminPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Test cases
	testCases := []GenericTestCase{
		{
			Payload: map[string]interface{}{
				"full_name": "Gerald Soto",
				"email":     "Not an email",
				"password":  "short",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Payload: map[string]interface{}{
				"full_name": "Gerald Soto",
				"email":     "gerald.soto@gmail.com",
				"password":  "gerald/password/2023",
			},
			ExpectedStatusCode: http.StatusCreated,
		},
		{
			Payload: map[string]interface{}{
				"full_name": "Gerald Soto",
				"email":     "gerald.soto@gmail.com",
				"password":  "gerald/password/2023",
			},
			ExpectedStatusCode: http.StatusConflict,
		},
	}

	// --- 1. Admin registers another admin ---
	for _, testCase := range testCases {
		w, r := PrepareRequest("POST", "/api/v1/accounts/admins", testCase.Payload)
		r.AddCookie(cookie)
		router.ServeHTTP(w, r)
		c.Equal(testCase.ExpectedStatusCode, w.Code)
	}

	// 2. --- Non-authenticated user tries to register an admin ---
	for _, testCase := range testCases {
		w, r := PrepareRequest("POST", "/api/v1/accounts/admins", testCase.Payload)
		router.ServeHTTP(w, r)
		c.Equal(http.StatusUnauthorized, w.Code)
	}

	// --- 3. Non-admin tries to register an admin ---
	w, r = PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredStudentEmail,
		"password": registeredStudentPass,
	})
	router.ServeHTTP(w, r)
	cookie = w.Result().Cookies()[0]

	for _, testCase := range testCases {
		w, r := PrepareRequest("POST", "/api/v1/accounts/admins", testCase.Payload)
		r.AddCookie(cookie)
		router.ServeHTTP(w, r)
		c.Equal(http.StatusForbidden, w.Code)
	}
}

func RegisterAdminAccount(req requests.RegisterAdminRequest) int {
	// Login as an admin
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredAdminEmail,
		"password": registeredAdminPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Register the new admin
	w, r = PrepareRequest("POST", "/api/v1/accounts/admins", map[string]interface{}{
		"full_name": req.FullName,
		"email":     req.Email,
		"password":  req.Password,
	})
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return w.Code
}

func TestRegisterTeacher(t *testing.T) {
	c := require.New(t)

	// Test cases
	testCases := []GenericTestCase{
		{
			Payload: map[string]interface{}{
				"full_name": "Zeeshan Glover",
				"email":     "Not an email",
				"password":  "short",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Payload: map[string]interface{}{
				"full_name": "Zeeshan Glover",
				"email":     "zeeshan.glover.2020@upb.edu.co",
				"password":  "zeeshan/password/2023",
			},
			ExpectedStatusCode: http.StatusCreated,
		},
		{
			Payload: map[string]interface{}{
				"full_name": "Zeeshan Glover",
				"email":     "zeeshan.glover.2020@upb.edu.co",
				"password":  "zeeshan/password/2023",
			},
			ExpectedStatusCode: http.StatusConflict,
		},
	}

	// Login as an admin
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredAdminEmail,
		"password": registeredAdminPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// --- 1. Admin registers a teacher ---
	for _, testCase := range testCases {
		w, r := PrepareRequest("POST", "/api/v1/accounts/teachers", testCase.Payload)
		r.AddCookie(cookie)
		router.ServeHTTP(w, r)
		c.Equal(testCase.ExpectedStatusCode, w.Code)
	}

	// 2. --- Non-authenticated user tries to register a teacher ---
	for _, testCase := range testCases {
		w, r := PrepareRequest("POST", "/api/v1/accounts/teachers", testCase.Payload)
		router.ServeHTTP(w, r)
		c.Equal(http.StatusUnauthorized, w.Code)
	}

	// --- 3. Non-admin tries to register a teacher ---

	// Login as a student
	w, r = PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredStudentEmail,
		"password": registeredStudentPass,
	})
	router.ServeHTTP(w, r)
	cookie = w.Result().Cookies()[0]

	for _, testCase := range testCases {
		w, r := PrepareRequest("POST", "/api/v1/accounts/teachers", testCase.Payload)
		r.AddCookie(cookie)
		router.ServeHTTP(w, r)
		c.Equal(http.StatusForbidden, w.Code)
	}
}

func RegisterTeacherAccount(req requests.RegisterTeacherRequest) int {
	// Login as an admin
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredAdminEmail,
		"password": registeredAdminPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Register the new teacher
	w, r = PrepareRequest("POST", "/api/v1/accounts/teachers", map[string]interface{}{
		"full_name": req.FullName,
		"email":     req.Email,
		"password":  req.Password,
	})
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)
	return w.Code
}

func TestListAdmins(t *testing.T) {
	c := require.New(t)

	// --- 1. Login as an admin ---
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredAdminEmail,
		"password": registeredAdminPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// List admins
	w, r = PrepareRequest("GET", "/api/v1/accounts/admins", nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)
	jsonResponse := ParseJsonResponse(w.Body)
	adminsList := jsonResponse["admins"].([]interface{})
	c.Equal(http.StatusOK, w.Code)
	c.GreaterOrEqual(len(adminsList), 1)

	for _, a := range adminsList {
		admin := a.(map[string]interface{})
		c.NotEmpty(admin["full_name"])
		c.NotEmpty(admin["created_at"])
	}

	// --- 2. Non-admin tries to list admins ---
	w, r = PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredStudentEmail,
		"password": registeredStudentPass,
	})
	router.ServeHTTP(w, r)
	cookie = w.Result().Cookies()[0]

	w, r = PrepareRequest("GET", "/api/v1/accounts/admins", nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)
	c.Equal(http.StatusForbidden, w.Code)
}
