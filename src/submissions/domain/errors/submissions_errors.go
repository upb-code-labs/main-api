package errors

import "net/http"

type StudentCannotSubmitToTestBlock struct{}

func (err StudentCannotSubmitToTestBlock) Error() string {
	return "You cannot submit to this test block"
}

func (err StudentCannotSubmitToTestBlock) StatusCode() int {
	return http.StatusForbidden
}

type UnableToQueueSubmissionWork struct{}

func (err UnableToQueueSubmissionWork) Error() string {
	return "Unable to queue your submission to be processed"
}

func (err UnableToQueueSubmissionWork) StatusCode() int {
	return http.StatusInternalServerError
}
