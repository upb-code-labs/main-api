package requests

import "github.com/UPB-Code-Labs/main-api/src/accounts/domain/dtos"

type RegisterAdminRequest struct {
	FullName string `json:"full_name" validate:"required,min=4,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255,secure_password"`
}

func (request *RegisterAdminRequest) ToDTO() *dtos.RegisterUserDTO {
	return &dtos.RegisterUserDTO{
		FullName: request.FullName,
		Email:    request.Email,
		Password: request.Password,
	}
}

type RegisterTeacherRequest struct {
	FullName string `json:"full_name" validate:"required,min=4,max=255"`
	Email    string `json:"email" validate:"required,email,institutional_email"`
	Password string `json:"password" validate:"required,min=8,max=255,secure_password"`
}

func (request *RegisterTeacherRequest) ToDTO() *dtos.RegisterUserDTO {
	return &dtos.RegisterUserDTO{
		FullName: request.FullName,
		Email:    request.Email,
		Password: request.Password,
	}
}

type RegisterUserRequest struct {
	FullName        string `json:"full_name" validate:"required,min=4,max=255"`
	Email           string `json:"email" validate:"required,email,institutional_email"`
	InstitutionalId string `json:"institutional_id" validate:"required,numeric,min=6,max=9"`
	Password        string `json:"password" validate:"required,min=8,max=255,secure_password"`
}

func (request *RegisterUserRequest) ToDTO() *dtos.RegisterUserDTO {
	return &dtos.RegisterUserDTO{
		FullName:        request.FullName,
		Email:           request.Email,
		InstitutionalId: request.InstitutionalId,
		Password:        request.Password,
	}
}
