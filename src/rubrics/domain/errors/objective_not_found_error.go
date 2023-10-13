package errors

import "net/http"

type ObjectiveNotFoundError struct{}

func (err *ObjectiveNotFoundError) Error() string {
	return "Rubric objective not found"
}

func (err *ObjectiveNotFoundError) StatusCode() int {
	return http.StatusNotFound
}
