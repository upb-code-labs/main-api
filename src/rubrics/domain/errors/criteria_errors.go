package errors

import "net/http"

type CriteriaNotFoundError struct{}

func (err *CriteriaNotFoundError) Error() string {
	return "Rubric criteria not found"
}

func (err *CriteriaNotFoundError) StatusCode() int {
	return http.StatusNotFound
}

// CriteriaDoesNotBelongToObjectiveError error to be thrown when teacher tries to set a criteria
// to a student's grade with a criteria that does not belong to the given objective
type CriteriaDoesNotBelongToObjectiveError struct{}

func (err *CriteriaDoesNotBelongToObjectiveError) Error() string {
	return "The criteria does not belong to the objective"
}

func (err *CriteriaDoesNotBelongToObjectiveError) StatusCode() int {
	return http.StatusBadRequest
}
