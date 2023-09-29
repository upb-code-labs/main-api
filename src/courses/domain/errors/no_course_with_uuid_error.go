package errors

import (
	"fmt"
	"net/http"
)

type NoCourseWithUUIDFound struct {
	UUID string
}

func (err NoCourseWithUUIDFound) Error() string {
	return fmt.Sprintf("Course with UUID %s not found", err.UUID)
}

func (err NoCourseWithUUIDFound) StatusCode() int {
	return http.StatusNotFound
}
