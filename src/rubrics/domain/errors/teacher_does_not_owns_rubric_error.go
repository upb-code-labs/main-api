package errors

import "net/http"

type TeacherDoesNotOwnsRubric struct{}

func (err *TeacherDoesNotOwnsRubric) Error() string {
	return "You do not own the rubric"
}

func (err *TeacherDoesNotOwnsRubric) StatusCode() int {
	return http.StatusForbidden
}
