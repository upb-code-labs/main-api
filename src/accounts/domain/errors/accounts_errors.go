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

type InstitutionalIdAlreadyInUseError struct {
	InstitutionalId string
}

func (err InstitutionalIdAlreadyInUseError) Error() string {
	return fmt.Sprintf("Institutional ID %s is already in use", err.InstitutionalId)
}

func (err InstitutionalIdAlreadyInUseError) StatusCode() int {
	return http.StatusConflict
}

type UserNotFoundError struct {
	Uuuid string
}

func (err UserNotFoundError) Error() string {
	return fmt.Sprintf("User with UUID %s not found", err.Uuuid)
}

func (err UserNotFoundError) StatusCode() int {
	return http.StatusNotFound
}
