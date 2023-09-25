package entities

type User struct {
	UUID            string
	Role            string
	FullName        string
	Email           string
	InstitutionalId string
	PasswordHash    string
}
