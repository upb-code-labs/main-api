package errors

import "net/http"

type StudentCannotSubmitToTestBlock struct{}

func (err StudentCannotSubmitToTestBlock) Error() string {
	return "You cannot submit to this test block"
}

func (err StudentCannotSubmitToTestBlock) StatusCode() int {
	return http.StatusForbidden
}
