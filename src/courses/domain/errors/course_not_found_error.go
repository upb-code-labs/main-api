package errors

import (
	"fmt"
	"net/http"
)

type CourseNotFoundError struct {
	UUID string
}

func (err CourseNotFoundError) Error() string {
	return fmt.Sprintf("Course with UUID %s not found", err.UUID)
}

func (err CourseNotFoundError) StatusCode() int {
	return http.StatusNotFound
}
