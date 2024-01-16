package errors

import "net/http"

type ObjectiveNotFoundError struct{}

func (err *ObjectiveNotFoundError) Error() string {
	return "Rubric objective not found"
}

func (err *ObjectiveNotFoundError) StatusCode() int {
	return http.StatusNotFound
}

// ObjectiveDoesNotBelongToRubricError error to be thrown when teacher tries to set a criteria
// to a student's grade with a criteria that belongs to an objective that does not belong to
// the current laboratory's rubric
type ObjectiveDoesNotBelongToRubricError struct{}

func (err *ObjectiveDoesNotBelongToRubricError) Error() string {
	return "The objective does not belong to the rubric"
}

func (err *ObjectiveDoesNotBelongToRubricError) StatusCode() int {
	return http.StatusBadRequest
}
