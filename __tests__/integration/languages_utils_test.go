package integration

import (
	"fmt"
	"net/http"
)

func GetSupportedLanguages(cookie *http.Cookie) (response map[string]interface{}, statusCode int) {
	w, r := PrepareRequest("GET", "/api/v1/languages", nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)
	return ParseJsonResponse(w.Body), w.Code
}

func GetLanguageTemplate(cookie *http.Cookie, uuid string) (bytes []byte, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/languages/%s/template", uuid)
	w, r := PrepareRequest("GET", endpoint, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)
	return w.Body.Bytes(), w.Code
}
