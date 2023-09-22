package entities

type User struct {
	UUID            string
	RoleUUID        string
	FullName        string
	Email           string
	InstitutionalId string
	PasswordHash    string
}
