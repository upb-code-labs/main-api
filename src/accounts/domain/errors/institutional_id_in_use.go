package errors

import (
	"fmt"
	"net/http"
)

type InstitutionalIdAlreadyInUseError struct {
	InstitutionalId string
}

func (err InstitutionalIdAlreadyInUseError) Error() string {
	return fmt.Sprintf("Institutional ID %s is already in use", err.InstitutionalId)
}

func (err InstitutionalIdAlreadyInUseError) StatusCode() int {
	return http.StatusConflict
}
