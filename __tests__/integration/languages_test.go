package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListSupportedLanguages(t *testing.T) {
	c := require.New(t)

	// Login as an student
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredStudentEmail,
		"password": registeredStudentPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Get the supported languages
	response, statusCode := GetSupportedLanguages(cookie)
	c.Equal(http.StatusOK, statusCode)

	// Check the response
	languages := response["languages"].([]interface{})
	c.Greater(len(languages), 0)

	// Get all the languages names
	var languagesNames []string
	for _, language := range languages {
		languagesNames = append(languagesNames, language.(map[string]interface{})["name"].(string))
	}

	// Check if the supported languages are included
	c.Contains(languagesNames, "Java")
}
