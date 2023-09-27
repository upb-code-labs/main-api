package errors

import (
	"net/http"
)

type UnauthorizedError struct {
	Message string
}

func (err UnauthorizedError) Error() string {
	return err.Message
}

func (err UnauthorizedError) StatusCode() int {
	return http.StatusUnauthorized
}
