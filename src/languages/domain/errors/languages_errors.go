package errors

import "net/http"

type LangNotFoundError struct{}

func (err *LangNotFoundError) Error() string {
	return "Language not found"
}

func (err *LangNotFoundError) StatusCode() int {
	return http.StatusNotFound
}
