package errors

import "net/http"

type TeacherDoesNotOwnBlock struct{}

func (err TeacherDoesNotOwnBlock) Error() string {
	return "You are not the owner of this block"
}

func (err TeacherDoesNotOwnBlock) StatusCode() int {
	return http.StatusForbidden
}

type BlockNotFound struct{}

func (err BlockNotFound) Error() string {
	return "No block was found with the given UUID"
}

func (err BlockNotFound) StatusCode() int {
	return http.StatusNotFound
}
