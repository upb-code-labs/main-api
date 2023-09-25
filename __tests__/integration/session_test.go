package integration

import (
	"testing"

	"github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure/requests"
	"github.com/stretchr/testify/require"
)

func TestLogin(t *testing.T) {
	c := require.New(t)

	// Register an student
	registerStudentPayload := requests.RegisterUserRequest{
		FullName:        "Delia Conn",
		Email:           "delia.conn.2020@upb.edu.co",
		InstitutionalId: "000149536",
		Password:        "delia/password/2023",
	}
	code := RegisterStudent(registerStudentPayload)
	c.Equal(201, code)

	// Register an admin
	registerAdminPayload := requests.RegisterAdminRequest{
		FullName: "Idun Yevhen",
		Email:    "idun.yevhen.2020@gmail.com",
		Password: "idun/password/2023",
	}
	code = RegisterAdminWithoutAuth(registerAdminPayload)
	c.Equal(201, code)

	// Login with an student
	w, r := PrepareRequest("POST", "/session/login", map[string]interface{}{
		"email":    registerStudentPayload.Email,
		"password": registerStudentPayload.Password,
	})
	router.ServeHTTP(w, r)
	jsonResponse := ParseJsonResponse(w.Body)
	responseUser := jsonResponse["user"].(map[string]interface{})

	c.Equal(200, w.Code)
	hasCookie := len(w.Result().Cookies()) == 1
	c.True(hasCookie)
	c.Equal(registerStudentPayload.FullName, responseUser["full_name"])
	c.Equal("student", responseUser["role"])

	// Login with an admin
	w, r = PrepareRequest("POST", "/session/login", map[string]interface{}{
		"email":    registerAdminPayload.Email,
		"password": registerAdminPayload.Password,
	})
	router.ServeHTTP(w, r)
	jsonResponse = ParseJsonResponse(w.Body)
	responseUser = jsonResponse["user"].(map[string]interface{})

	c.Equal(200, w.Code)
	hasCookie = len(w.Result().Cookies()) == 1
	c.True(hasCookie)
	c.Equal(registerAdminPayload.FullName, responseUser["full_name"])
	c.Equal("admin", responseUser["role"])

	// Login with wrong credentials
	w, r = PrepareRequest("POST", "/session/login", map[string]interface{}{
		"email":    registerAdminPayload.Email,
		"password": "wrong password",
	})
	router.ServeHTTP(w, r)

	c.Equal(401, w.Code)
	hasCookie = len(w.Result().Cookies()) == 1
	c.False(hasCookie)
}
