package integration

import (
	"net/http"
	"testing"

	"github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure/requests"
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

func TestGetCreatedRubrics(t *testing.T) {
	c := require.New(t)

	// Register a teacher
	testTeacherEmail := "nirmala.ivona.2020@upb.edu.co"
	testTeacherPass := "nirmala/password/2020"
	code := RegisterTeacherAccount(requests.RegisterTeacherRequest{
		FullName: "Nirmala Ivona",
		Email:    testTeacherEmail,
		Password: testTeacherPass,
	})
	c.Equal(201, code)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    testTeacherEmail,
		"password": testTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Get created rubrics
	response, status := GetCreatedRubrics(cookie)
	rubrics := response["rubrics"].([]interface{})
	c.Equal(http.StatusOK, status)
	c.Equal(0, len(rubrics))
	c.NotEmpty(response["message"])

	// Create a rubric
	_, status = CreateRubric(cookie, map[string]interface{}{
		"name": "Rubric 1",
	})
	c.Equal(http.StatusCreated, status)

	// Get created rubrics
	response, status = GetCreatedRubrics(cookie)
	rubrics = response["rubrics"].([]interface{})
	c.Equal(http.StatusOK, status)
	c.Equal(1, len(rubrics))

	// Validate rubric fields
	rubric := rubrics[0].(map[string]interface{})
	c.NotEmpty(rubric["uuid"])
	c.NotEmpty(rubric["name"])
}

func GetCreatedRubrics(cookie *http.Cookie) (response map[string]interface{}, status int) {
	w, r := PrepareRequest("GET", "/api/v1/rubrics", nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return ParseJsonResponse(w.Body), w.Code
}

func TestGetRubricByUUID(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a rubric
	response, status := CreateRubric(cookie, map[string]interface{}{
		"name": "Rubric 1",
	})
	c.Equal(http.StatusCreated, status)
	rubricUUID := response["uuid"].(string)

	// Create a teacher
	testTeacherEmail := "henriette.otylia.2020@upb.edu.co"
	testTeacherPass := "henriette/password/2020"
	code := RegisterTeacherAccount(requests.RegisterTeacherRequest{
		FullName: "Henriette Otylia",
		Email:    testTeacherEmail,
		Password: testTeacherPass,
	})
	c.Equal(201, code)

	// Login as the new teacher
	w, r = PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    testTeacherEmail,
		"password": testTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie = w.Result().Cookies()[0]

	// Create a rubric
	response, status = CreateRubric(cookie, map[string]interface{}{
		"name": "Rubric 2",
	})
	c.Equal(http.StatusCreated, status)
	rubricUUID2 := response["uuid"].(string)

	// Test cases
	testCases := []GenericTestCase{
		GenericTestCase{
			Payload: map[string]interface{}{
				"rubricUUID": "not-valid-uuid",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				"rubricUUID": "90b2edf3-72fc-4682-be1c-1274c70785d9",
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				"rubricUUID": rubricUUID,
			},
			ExpectedStatusCode: http.StatusForbidden,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				"rubricUUID": rubricUUID2,
			},
			ExpectedStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		response, status := GetRubricByUUID(cookie, testCase.Payload["rubricUUID"].(string))
		c.Equal(testCase.ExpectedStatusCode, status)

		if testCase.ExpectedStatusCode == http.StatusOK {
			rubric := response["rubric"].(map[string]interface{})
			c.NotEmpty(response["message"])
			c.Equal(rubricUUID2, rubric["uuid"])
			c.Equal("Rubric 2", rubric["name"])

			objective := rubric["objectives"].([]interface{})[0].(map[string]interface{})
			c.NotEmpty(objective["uuid"])
			c.NotEmpty(objective["description"])

			criteria := objective["criteria"].([]interface{})[0].(map[string]interface{})
			c.NotEmpty(criteria["uuid"])
			c.NotEmpty(criteria["description"])
			c.NotEmpty(criteria["weight"])
		}
	}
}

func GetRubricByUUID(cookie *http.Cookie, uuid string) (response map[string]interface{}, status int) {
	w, r := PrepareRequest("GET", "/api/v1/rubrics/"+uuid, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return ParseJsonResponse(w.Body), w.Code
}
