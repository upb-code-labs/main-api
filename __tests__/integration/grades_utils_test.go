package integration

import (
	"fmt"
	"net/http"
)

func GetSummarizedGrades(laboratoryUUID string, cookie *http.Cookie) (response map[string]interface{}, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/grades/laboratories/%s", laboratoryUUID)
	w, r := PrepareRequest("GET", endpoint, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

type SetCriteriaToStudentGradeUtilsDTO struct {
	LaboratoryUUID string
	StudentUUID    string
	ObjectiveUUID  string
	CriteriaUUID   string
}

func SetCriteriaToStudentGrade(dto *SetCriteriaToStudentGradeUtilsDTO, cookie *http.Cookie) (response map[string]interface{}, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/grades/laboratories/%s/students/%s", dto.LaboratoryUUID, dto.StudentUUID)
	w, r := PrepareRequest("PUT", endpoint, map[string]interface{}{
		"objective_uuid": dto.ObjectiveUUID,
		"criteria_uuid":  dto.CriteriaUUID,
	})
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

type GetStudentGradeUtilsDTO struct {
	LaboratoryUUID string
	StudentUUID    string
	RubricUUID     string
}

func GetStudentGrade(dto *GetStudentGradeUtilsDTO, cookie *http.Cookie) (response map[string]interface{}, statusCode int) {
	endpoint := fmt.Sprintf(
		"/api/v1/grades/laboratories/%s/students/%s/rubrics/%s",
		dto.LaboratoryUUID,
		dto.StudentUUID,
		dto.RubricUUID,
	)

	w, r := PrepareRequest("GET", endpoint, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

type SetCommentToStudentGradeUtilsDTO struct {
	LaboratoryUUID string
	StudentUUID    string
	Comment        string
}

func SetCommentToStudentGrade(dto *SetCommentToStudentGradeUtilsDTO, cookie *http.Cookie) (response map[string]interface{}, statusCode int) {
	endpoint := fmt.Sprintf(
		"/api/v1/grades/laboratories/%s/students/%s/comment",
		dto.LaboratoryUUID,
		dto.StudentUUID,
	)

	w, r := PrepareRequest("PUT", endpoint, map[string]interface{}{
		"comment": dto.Comment,
	})
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}
