package errors

import "net/http"

type CriteriaNotFoundError struct{}

func (err *CriteriaNotFoundError) Error() string {
	return "Rubric criteria not found"
}

func (err *CriteriaNotFoundError) StatusCode() int {
	return http.StatusNotFound
}
