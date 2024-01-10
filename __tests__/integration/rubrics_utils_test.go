package integration

import "net/http"

func CreateRubric(cookie *http.Cookie, payload map[string]interface{}) (response map[string]interface{}, status int) {
	w, r := PrepareRequest("POST", "/api/v1/rubrics", payload)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return ParseJsonResponse(w.Body), w.Code
}

func GetRubricsCreatedByUser(cookie *http.Cookie) (response map[string]interface{}, status int) {
	w, r := PrepareRequest("GET", "/api/v1/rubrics", nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return ParseJsonResponse(w.Body), w.Code
}

func GetRubricByUUID(cookie *http.Cookie, uuid string) (response map[string]interface{}, status int) {
	w, r := PrepareRequest("GET", "/api/v1/rubrics/"+uuid, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return ParseJsonResponse(w.Body), w.Code
}

func DeleteRubric(cookie *http.Cookie, uuid string) (response map[string]interface{}, status int) {
	endpoint := "/api/v1/rubrics/" + uuid
	w, r := PrepareRequest("DELETE", endpoint, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return ParseJsonResponse(w.Body), w.Code
}

func UpdateRubricName(cookie *http.Cookie, uuid string, payload map[string]interface{}) (response map[string]interface{}, status int) {
	w, r := PrepareRequest("PATCH", "/api/v1/rubrics/"+uuid+"/name", payload)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return ParseJsonResponse(w.Body), w.Code

}

func AddObjectiveToRubric(cookie *http.Cookie, rubricUUID string, payload map[string]interface{}) (response map[string]interface{}, status int) {
	w, r := PrepareRequest("POST", "/api/v1/rubrics/"+rubricUUID+"/objectives", payload)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return ParseJsonResponse(w.Body), w.Code
}

func UpdateObjective(cookie *http.Cookie, objectiveUUID string, payload map[string]interface{}) (response map[string]interface{}, status int) {
	w, r := PrepareRequest("PUT", "/api/v1/rubrics/objectives/"+objectiveUUID, payload)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return ParseJsonResponse(w.Body), w.Code
}

func DeleteObjective(cookie *http.Cookie, objectiveUUID string) (response map[string]interface{}, status int) {
	w, r := PrepareRequest("DELETE", "/api/v1/rubrics/objectives/"+objectiveUUID, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return ParseJsonResponse(w.Body), w.Code
}

func AddCriteriaToObjective(cookie *http.Cookie, objectiveUUID string, payload map[string]interface{}) (response map[string]interface{}, status int) {
	w, r := PrepareRequest("POST", "/api/v1/rubrics/objectives/"+objectiveUUID+"/criteria", payload)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return ParseJsonResponse(w.Body), w.Code
}

func UpdateCriteria(cookie *http.Cookie, criteriaUUID string, payload map[string]interface{}) (response map[string]interface{}, status int) {
	w, r := PrepareRequest("PUT", "/api/v1/rubrics/criteria/"+criteriaUUID, payload)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return ParseJsonResponse(w.Body), w.Code
}

func DeleteCriteria(cookie *http.Cookie, criteriaUUID string) (response map[string]interface{}, status int) {
	w, r := PrepareRequest("DELETE", "/api/v1/rubrics/criteria/"+criteriaUUID, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	return ParseJsonResponse(w.Body), w.Code
}
