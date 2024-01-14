package dtos

type RegisterUserDTO struct {
	FullName        string
	Email           string
	InstitutionalId string
	Password        string
	CreatedBy       string
}

type UpdatePasswordDTO struct {
	UserUUID    string
	OldPassword string
	NewPassword string
}
