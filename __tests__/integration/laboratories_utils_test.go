package integration

import "net/http"

func CreateLaboratory(cookie *http.Cookie, payload map[string]interface{}) (response map[string]interface{}, statusCode int) {
	w, r := PrepareRequest("POST", "/api/v1/laboratories", payload)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func GetLaboratoryByUUID(cookie *http.Cookie, uuid string) (response map[string]interface{}, statusCode int) {
	w, r := PrepareRequest("GET", "/api/v1/laboratories/"+uuid, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func UpdateLaboratory(cookie *http.Cookie, uuid string, payload map[string]interface{}) (response map[string]interface{}, statusCode int) {
	w, r := PrepareRequest("PUT", "/api/v1/laboratories/"+uuid, payload)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func CreateMarkdownBlock(cookie *http.Cookie, laboratoryUUID string) (response map[string]interface{}, statusCode int) {
	w, r := PrepareRequest("POST", "/api/v1/laboratories/markdown_blocks/"+laboratoryUUID, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}
