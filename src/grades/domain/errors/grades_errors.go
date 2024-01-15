package errors

import "net/http"

// LaboratoryDoesNotHaveRubricError error to be thrown when a laboratory does not have a rubric
type LaboratoryDoesNotHaveRubricError struct{}

func (err LaboratoryDoesNotHaveRubricError) Error() string {
	return "The laboratory does not have a rubric"
}

func (err LaboratoryDoesNotHaveRubricError) StatusCode() int {
	return http.StatusBadRequest
}

// RubricDoesNotMatchLaboratoryError error to be thrown when teacher tries to set a criteria
// to a student's grade with a rubric that does not match the laboratory's rubric
type RubricDoesNotMatchLaboratoryError struct{}

func (err RubricDoesNotMatchLaboratoryError) Error() string {
	return "The rubric you are trying to use to grade the laboratory does not match the laboratory's rubric"
}

func (err RubricDoesNotMatchLaboratoryError) StatusCode() int {
	return http.StatusBadRequest
}
