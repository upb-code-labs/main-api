package errors

import (
	"net/http"
)

type LaboratoryNotFoundError struct{}

func (err LaboratoryNotFoundError) Error() string {
	return "No laboratory found with the given id"
}

func (err LaboratoryNotFoundError) StatusCode() int {
	return http.StatusNotFound
}
