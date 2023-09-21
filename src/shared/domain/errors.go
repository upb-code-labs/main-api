package domain

type DomainError interface {
	Error() string
	StatusCode() int
}
