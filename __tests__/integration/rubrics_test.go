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
	response, status := GetRubricsCreatedByUser(cookie)
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
	response, status = GetRubricsCreatedByUser(cookie)
	rubrics = response["rubrics"].([]interface{})
	c.Equal(http.StatusOK, status)
	c.Equal(1, len(rubrics))

	// Validate rubric fields
	rubric := rubrics[0].(map[string]interface{})
	c.NotEmpty(rubric["uuid"])
	c.NotEmpty(rubric["name"])
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

func TestUpdateRubricName(t *testing.T) {
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

	// Test cases
	newName := "New name"
	testCases := []GenericTestCase{
		GenericTestCase{
			Payload: map[string]interface{}{
				"name": "a",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				"name": newName,
			},
			ExpectedStatusCode: http.StatusNoContent,
		},
	}

	for _, testCase := range testCases {
		_, status := UpdateRubricName(cookie, rubricUUID, testCase.Payload)
		c.Equal(testCase.ExpectedStatusCode, status)
	}

	// Get rubric
	response, status = GetRubricByUUID(cookie, rubricUUID)
	c.Equal(http.StatusOK, status)

	rubric := response["rubric"].(map[string]interface{})
	c.Equal(newName, rubric["name"])
	c.NotEmpty(rubric["uuid"])
}

func TestAddObjectiveToRubric(t *testing.T) {
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

	// Test cases
	objectiveDescription := "Objective 1"
	testCases := []GenericTestCase{
		GenericTestCase{
			Payload: map[string]interface{}{
				"rubricUUID":  "not-valid-uuid",
				"description": "Objective 1",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				"rubricUUID":  rubricUUID,
				"description": "short",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				"rubricUUID":  rubricUUID,
				"description": objectiveDescription,
			},
			ExpectedStatusCode: http.StatusCreated,
		},
	}

	for _, testCase := range testCases {
		response, status := AddObjectiveToRubric(cookie, testCase.Payload["rubricUUID"].(string), testCase.Payload)
		c.Equal(testCase.ExpectedStatusCode, status)

		if testCase.ExpectedStatusCode == http.StatusCreated {
			c.NotEmpty(response["uuid"])
			c.NotEmpty(response["message"])
		}
	}

	// Get rubric
	response, status = GetRubricByUUID(cookie, rubricUUID)
	c.Equal(http.StatusOK, status)

	rubric := response["rubric"].(map[string]interface{})
	c.Equal(2, len(rubric["objectives"].([]interface{})))

	objective := rubric["objectives"].([]interface{})[1].(map[string]interface{})
	c.Equal(objectiveDescription, objective["description"])
	c.NotEmpty(objective["uuid"])
	c.Empty(objective["criteria"])
}

func TestUpdateObjective(t *testing.T) {
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

	// Create an objective
	response, status = AddObjectiveToRubric(cookie, rubricUUID, map[string]interface{}{
		"description": "Old description",
	})
	c.Equal(http.StatusCreated, status)
	objectiveUUID := response["uuid"].(string)

	// Test cases
	newDescription := "New description"
	testCases := []GenericTestCase{
		GenericTestCase{
			Payload: map[string]interface{}{
				"description": "short",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				"description": newDescription,
			},
			ExpectedStatusCode: http.StatusNoContent,
		},
	}

	for _, testCase := range testCases {
		_, status := UpdateObjective(cookie, objectiveUUID, testCase.Payload)
		c.Equal(testCase.ExpectedStatusCode, status)
	}

	// Get rubric
	response, status = GetRubricByUUID(cookie, rubricUUID)
	c.Equal(http.StatusOK, status)

	rubric := response["rubric"].(map[string]interface{})
	c.Equal(2, len(rubric["objectives"].([]interface{})))

	objective := rubric["objectives"].([]interface{})[1].(map[string]interface{})
	c.Equal(newDescription, objective["description"])
	c.NotEmpty(objective["uuid"])
}

func TestDeleteObjective(t *testing.T) {
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

	// Get the rubric
	response, status = GetRubricByUUID(cookie, rubricUUID)
	c.Equal(http.StatusOK, status)

	rubric := response["rubric"].(map[string]interface{})
	c.Equal(1, len(rubric["objectives"].([]interface{})))

	objective := rubric["objectives"].([]interface{})[0].(map[string]interface{})
	objectiveUUID := objective["uuid"].(string)

	// Test cases
	testCases := []GenericTestCase{
		GenericTestCase{
			Payload: map[string]interface{}{
				"objectiveUUID": "not-valid-uuid",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				"objectiveUUID": "8bf7e0a1-6475-41b6-81f6-e804484a9d67",
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				"objectiveUUID": objectiveUUID,
			},
			ExpectedStatusCode: http.StatusNoContent,
		},
	}

	for _, testCase := range testCases {
		_, status := DeleteObjective(cookie, testCase.Payload["objectiveUUID"].(string))
		c.Equal(testCase.ExpectedStatusCode, status)
	}

	// Get rubric
	response, status = GetRubricByUUID(cookie, rubricUUID)
	c.Equal(http.StatusOK, status)

	rubric = response["rubric"].(map[string]interface{})
	c.Equal(0, len(rubric["objectives"].([]interface{})))
}

func TestAddCriteriaToObjective(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	firstTeacherCookie := w.Result().Cookies()[0]

	// Create a rubric
	response, status := CreateRubric(firstTeacherCookie, map[string]interface{}{
		"name": "Rubric 1",
	})
	c.Equal(http.StatusCreated, status)
	firstTeacherRubricUUID := response["uuid"].(string)

	// Create an objective
	response, status = AddObjectiveToRubric(firstTeacherCookie, firstTeacherRubricUUID, map[string]interface{}{
		"description": "Objective 1",
	})
	c.Equal(http.StatusCreated, status)
	firstTeacherObjectiveUUID := response["uuid"].(string)

	// Login as the second teacher
	w, r = PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    secondRegisteredTeacherEmail,
		"password": secondRegisteredTeacherPass,
	})
	router.ServeHTTP(w, r)
	secondTeacherCookie := w.Result().Cookies()[0]

	// Create a rubric
	response, status = CreateRubric(secondTeacherCookie, map[string]interface{}{
		"name": "Rubric 2",
	})
	c.Equal(http.StatusCreated, status)
	secondTeacherRubricUUID := response["uuid"].(string)

	// Create an objective
	response, status = AddObjectiveToRubric(secondTeacherCookie, secondTeacherRubricUUID, map[string]interface{}{
		"description": "Objective 2",
	})
	c.Equal(http.StatusCreated, status)
	secondTeacherObjectiveUUID := response["uuid"].(string)

	// Test cases
	testCases := []GenericTestCase{
		GenericTestCase{
			Payload: map[string]interface{}{
				"objectiveUUID": "not-valid-uuid",
				"description":   "Criteria 1",
				"weight":        5.00,
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				"objectiveUUID": firstTeacherObjectiveUUID,
				"description":   "Criteria 1",
				"weight":        5.00,
			},
			ExpectedStatusCode: http.StatusForbidden,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				"objectiveUUID": secondTeacherObjectiveUUID,
				"description":   "short",
				"weight":        5.00,
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				"objectiveUUID": "adc73ae3-80ad-45d3-ae23-bd81e6e0b805",
				"description":   "Criteria 1",
				"weight":        5.00,
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				"objectiveUUID": secondTeacherObjectiveUUID,
				"description":   "Criteria 1",
				"weight":        5.00,
			},
			ExpectedStatusCode: http.StatusCreated,
		},
	}

	for _, testCase := range testCases {
		response, status := AddCriteriaToObjective(secondTeacherCookie, testCase.Payload["objectiveUUID"].(string), testCase.Payload)
		c.Equal(testCase.ExpectedStatusCode, status)

		if testCase.ExpectedStatusCode == http.StatusCreated {
			c.NotEmpty(response["uuid"])
			c.NotEmpty(response["message"])
		}
	}

	// Get rubric
	response, status = GetRubricByUUID(secondTeacherCookie, secondTeacherRubricUUID)
	c.Equal(http.StatusOK, status)

	rubric := response["rubric"].(map[string]interface{})
	c.Equal(2, len(rubric["objectives"].([]interface{})))

	objective := rubric["objectives"].([]interface{})[1].(map[string]interface{})
	c.Equal(1, len(objective["criteria"].([]interface{})))

	criteria := objective["criteria"].([]interface{})[0].(map[string]interface{})
	c.Equal("Criteria 1", criteria["description"])
	c.NotEmpty(criteria["uuid"])
	c.NotEmpty(criteria["weight"])
}

func TestUpdateCriteria(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	firstTeacherCookie := w.Result().Cookies()[0]

	// Create a rubric
	response, status := CreateRubric(firstTeacherCookie, map[string]interface{}{
		"name": "Rubric 1",
	})
	c.Equal(http.StatusCreated, status)
	firstTeacherRubricUUID := response["uuid"].(string)

	// Get the criteria UUID
	response, status = GetRubricByUUID(firstTeacherCookie, firstTeacherRubricUUID)
	c.Equal(http.StatusOK, status)

	rubric := response["rubric"].(map[string]interface{})
	c.Equal(1, len(rubric["objectives"].([]interface{})))

	objective := rubric["objectives"].([]interface{})[0].(map[string]interface{})
	c.Equal(1, len(objective["criteria"].([]interface{})))

	criteria := objective["criteria"].([]interface{})[0].(map[string]interface{})
	criteriaUUID := criteria["uuid"].(string)

	// Test cases
	newDescription := "New description"
	newWeight := 0.125

	testCases := []GenericTestCase{
		GenericTestCase{
			Payload: map[string]interface{}{
				"description": "short",
				"weight":      0.125,
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				"description": newDescription,
				"weight":      newWeight,
			},
			ExpectedStatusCode: http.StatusNoContent,
		},
	}

	for _, testCase := range testCases {
		_, status := UpdateCriteria(firstTeacherCookie, criteriaUUID, testCase.Payload)
		c.Equal(testCase.ExpectedStatusCode, status)
	}
}

func TestDeleteCriteria(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	firstTeacherCookie := w.Result().Cookies()[0]

	// Create a rubric
	response, status := CreateRubric(firstTeacherCookie, map[string]interface{}{
		"name": "Rubric 1",
	})
	c.Equal(http.StatusCreated, status)
	firstTeacherRubricUUID := response["uuid"].(string)

	// Get the criteria UUID
	response, status = GetRubricByUUID(firstTeacherCookie, firstTeacherRubricUUID)
	c.Equal(http.StatusOK, status)

	rubric := response["rubric"].(map[string]interface{})
	c.Equal(1, len(rubric["objectives"].([]interface{})))

	objective := rubric["objectives"].([]interface{})[0].(map[string]interface{})
	c.Equal(1, len(objective["criteria"].([]interface{})))

	criteria := objective["criteria"].([]interface{})[0].(map[string]interface{})
	criteriaUUID := criteria["uuid"].(string)

	// Test cases
	testCases := []GenericTestCase{
		GenericTestCase{
			Payload: map[string]interface{}{
				"criteriaUUID": "not-valid-uuid",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				"criteriaUUID": "ceb714e2-3e97-4b25-a8b7-6f5e919a3e3c",
			},
			ExpectedStatusCode: http.StatusNotFound,
		},
		GenericTestCase{
			Payload: map[string]interface{}{
				"criteriaUUID": criteriaUUID,
			},
			ExpectedStatusCode: http.StatusNoContent,
		},
	}

	for _, testCase := range testCases {
		_, status := DeleteCriteria(firstTeacherCookie, testCase.Payload["criteriaUUID"].(string))
		c.Equal(testCase.ExpectedStatusCode, status)
	}

	// Get rubric
	response, status = GetRubricByUUID(firstTeacherCookie, firstTeacherRubricUUID)
	c.Equal(http.StatusOK, status)

	rubric = response["rubric"].(map[string]interface{})
	c.Equal(1, len(rubric["objectives"].([]interface{})))

	objective = rubric["objectives"].([]interface{})[0].(map[string]interface{})
	c.Equal(0, len(objective["criteria"].([]interface{})))
}
