package errors

import (
	"fmt"
	"net/http"
)

type UnauthorizedError struct {
	Message string
}

func (err UnauthorizedError) Error() string {
	return fmt.Sprintf("Unauthorized: %s", err.Message)
}

func (err UnauthorizedError) StatusCode() int {
	return http.StatusUnauthorized
}
