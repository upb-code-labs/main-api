package errors

import (
	"fmt"
	"net/http"
)

type NotEnoughPermissionsError struct {
	Message string
}

func (err NotEnoughPermissionsError) Error() string {
	return fmt.Sprintf("Not enough permissions: %s", err.Message)
}

func (err NotEnoughPermissionsError) StatusCode() int {
	return http.StatusForbidden
}
