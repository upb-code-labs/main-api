package errors

import (
	"net/http"
)

type UnchangedCourseNameError struct {
}

func (err UnchangedCourseNameError) Error() string {
	return "The course has the same name"
}

func (err UnchangedCourseNameError) StatusCode() int {
	return http.StatusBadRequest
}
