package errors

import (
	"fmt"
	"net/http"
)

type UserNotFoundError struct {
	Uuuid string
}

func (err UserNotFoundError) Error() string {
	return fmt.Sprintf("User with UUID %s not found", err.Uuuid)
}

func (err UserNotFoundError) StatusCode() int {
	return http.StatusNotFound
}
