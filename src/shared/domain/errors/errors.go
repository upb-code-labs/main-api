package errors

type DomainError interface {
	Error() string
	StatusCode() int
}
