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

type TeacherDoesNotOwnLaboratoryError struct{}

func (err TeacherDoesNotOwnLaboratoryError) Error() string {
	return "You don't own the laboratory"
}

func (err TeacherDoesNotOwnLaboratoryError) StatusCode() int {
	return http.StatusForbidden
}

type UserCannotAccessProgressSummaryError struct{}

func (err UserCannotAccessProgressSummaryError) Error() string {
	return "You cannot access the progress summary of the given student in the given laboratory"
}

func (err UserCannotAccessProgressSummaryError) StatusCode() int {
	return http.StatusForbidden
}
