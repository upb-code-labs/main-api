package integration

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure/requests"
)

func RegisterStudentAccount(req requests.RegisterUserRequest) int {
	w, r := PrepareRequest("POST", "/api/v1/accounts/students", map[string]interface{}{
		"full_name":        req.FullName,
		"email":            req.Email,
		"institutional_id": req.InstitutionalId,
		"password":         req.Password,
	})

	router.ServeHTTP(w, r)
	return w.Code
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

func SearchStudentsByFullName(cookie *http.Cookie, fullName string) (response map[string]interface{}, statusCode int) {
	w, r := PrepareRequest("GET", "/api/v1/accounts/students?fullName="+fullName, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)
	return ParseJsonResponse(w.Body), w.Code
}
