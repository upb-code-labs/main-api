package integration

import (
	"net/http"
	"testing"

	"github.com/gabriel-vasile/mimetype"
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
	c.Contains(languagesNames, "Java JDK 17")
}

func GetFirstSupportedLanguage(cookie *http.Cookie) (language map[string]interface{}) {
	w, r := PrepareRequest("GET", "/api/v1/languages", nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)

	languages := jsonResponse["languages"].([]interface{})
	language = languages[0].(map[string]interface{})
	return language
}

func TestGetLanguageTemplate(t *testing.T) {
	c := require.New(t)

	// Login as an student
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredStudentEmail,
		"password": registeredStudentPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Get a supported language from the supported languages list
	language := GetFirstSupportedLanguage(cookie)

	// Get the language template
	template, statusCode := GetLanguageTemplate(cookie, language["uuid"].(string))
	c.Equal(http.StatusOK, statusCode)

	// Check the response
	c.Greater(len(template), 0)

	// Check the MIMETYPE
	mtype := mimetype.Detect(template)
	c.Equal("application/zip", mtype.String())
}
