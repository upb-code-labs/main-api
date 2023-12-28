package errors

import "net/http"

type LangNotFoundError struct{}

func (err *LangNotFoundError) Error() string {
	return "Language not found"
}

func (err *LangNotFoundError) StatusCode() int {
	return http.StatusNotFound
}

type StaticFilesMicroserviceError struct {
	Code    int
	Message string
}

func (err *StaticFilesMicroserviceError) Error() string {
	return err.Message
}

func (err *StaticFilesMicroserviceError) StatusCode() int {
	return err.Code
}
