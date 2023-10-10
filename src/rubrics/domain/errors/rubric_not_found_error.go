package errors

import "net/http"

type RubricNotFoundError struct{}

func (err *RubricNotFoundError) Error() string {
	return "Rubric not found"
}

func (err *RubricNotFoundError) StatusCode() int {
	return http.StatusNotFound
}
