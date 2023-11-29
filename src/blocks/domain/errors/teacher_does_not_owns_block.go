package errors

import "net/http"

type TeacherDoesNotOwnBlock struct{}

func (err TeacherDoesNotOwnBlock) Error() string {
	return "You are not the owner of this block"
}

func (err TeacherDoesNotOwnBlock) StatusCode() int {
	return http.StatusForbidden
}
