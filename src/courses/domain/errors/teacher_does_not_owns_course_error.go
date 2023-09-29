package errors

import (
	"net/http"
)

type TeacherDoesNotOwnsCourseError struct {
}

func (err TeacherDoesNotOwnsCourseError) Error() string {
	return "You do not own the course"
}

func (err TeacherDoesNotOwnsCourseError) StatusCode() int {
	return http.StatusForbidden
}
