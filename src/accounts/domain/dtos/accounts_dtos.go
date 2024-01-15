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

type UpdateAccountDTO struct {
	UserUUID        string
	FullName        string
	Email           string
	InstitutionalId *string
	Password        string
}

type UserProfileDTO struct {
	FullName        string  `json:"full_name"`
	Email           string  `json:"email"`
	InstitutionalId *string `json:"institutional_id"`
}
