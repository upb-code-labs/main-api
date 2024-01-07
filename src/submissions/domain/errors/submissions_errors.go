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

type StudentSubmissionNotFound struct{}

func (err StudentSubmissionNotFound) Error() string {
	return "You don't have a submission for this test block"
}

func (err StudentSubmissionNotFound) StatusCode() int {
	return http.StatusNotFound
}

type StudentHasRecentSubmission struct{}

func (err StudentHasRecentSubmission) Error() string {
	return "You need to wait, at least, 1 minute before submitting again"
}

func (err StudentHasRecentSubmission) StatusCode() int {
	return http.StatusForbidden
}

type StudentHasPendingSubmission struct{}

func (err StudentHasPendingSubmission) Error() string {
	return "You already have a pending or running submission for this test block"
}

func (err StudentHasPendingSubmission) StatusCode() int {
	return http.StatusForbidden
}

type LaboratoryIsClosed struct{}

func (err LaboratoryIsClosed) Error() string {
	return "The laboratory no longer accepts submissions"
}

func (err LaboratoryIsClosed) StatusCode() int {
	return http.StatusForbidden
}
