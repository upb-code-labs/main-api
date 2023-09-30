package errors

import (
	"net/http"
)

type NotEnoughPermissionsError struct {
}

func (err NotEnoughPermissionsError) Error() string {
	return "Not enough permissions"
}

func (err NotEnoughPermissionsError) StatusCode() int {
	return http.StatusForbidden
}
