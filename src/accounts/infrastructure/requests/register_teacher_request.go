package requests

import "github.com/UPB-Code-Labs/main-api/src/accounts/domain/dtos"

type RegisterTeacherRequest struct {
	FullName string `json:"full_name" validate:"required,min=4,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=255,secure_password"`
}

func (request *RegisterTeacherRequest) ToDTO() *dtos.RegisterUserDTO {
	return &dtos.RegisterUserDTO{
		FullName: request.FullName,
		Email:    request.Email,
		Password: request.Password,
	}
}
