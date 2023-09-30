package errors

import (
	"net/http"
)

type UserNotInCourseError struct{}

func (err UserNotInCourseError) Error() string {
	return "You are not enrolled in the course"
}

func (err UserNotInCourseError) StatusCode() int {
	return http.StatusForbidden
}
