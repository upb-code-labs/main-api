package errors

import "net/http"

type TeacherDoesNotOwnsRubric struct{}

func (err *TeacherDoesNotOwnsRubric) Error() string {
	return "You do not own the rubric"
}

func (err *TeacherDoesNotOwnsRubric) StatusCode() int {
	return http.StatusForbidden
}

type RubricNotFoundError struct{}

func (err *RubricNotFoundError) Error() string {
	return "Rubric not found"
}

func (err *RubricNotFoundError) StatusCode() int {
	return http.StatusNotFound
}
