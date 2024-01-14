package errors

import (
	"net/http"
)

type InvalidCredentialsError struct{}

func (err InvalidCredentialsError) Error() string {
	return "Credentials are wrong"
}

func (err InvalidCredentialsError) StatusCode() int {
	return http.StatusUnauthorized
}
