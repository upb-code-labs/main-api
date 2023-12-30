package errors

type DomainError interface {
	Error() string
	StatusCode() int
}

type GenericDomainError struct {
	Code    int
	Message string
}

func (err *GenericDomainError) Error() string {
	return err.Message
}

func (err *GenericDomainError) StatusCode() int {
	return err.Code
}
