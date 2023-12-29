package errors

type DomainError interface {
	Error() string
	StatusCode() int
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
