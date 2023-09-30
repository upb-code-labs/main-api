package errors

import (
	"fmt"
	"net/http"
)

type NoCourseWithInvitationCodeError struct {
	Code string
}

func (err NoCourseWithInvitationCodeError) Error() string {
	return fmt.Sprintf("Course with invitation code %s not found", err.Code)
}

func (err NoCourseWithInvitationCodeError) StatusCode() int {
	return http.StatusNotFound
}
