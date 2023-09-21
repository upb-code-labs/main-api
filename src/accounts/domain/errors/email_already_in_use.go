package errors

import (
	"fmt"
	"net/http"
)

type EmailAlreadyInUseError struct {
	Email string
}

func (err EmailAlreadyInUseError) Error() string {
	return fmt.Sprintf("Email %s is already in use", err.Email)
}

func (err EmailAlreadyInUseError) StatusCode() int {
	return http.StatusConflict
}
