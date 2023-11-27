package integration

import "net/http"

func CreateLaboratory(cookie *http.Cookie, payload map[string]interface{}) (response map[string]interface{}, statusCode int) {
	w, r := PrepareRequest("POST", "/api/v1/laboratories", payload)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}
