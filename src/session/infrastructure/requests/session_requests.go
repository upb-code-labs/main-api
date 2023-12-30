package requests

import "github.com/UPB-Code-Labs/main-api/src/session/domain/dtos"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (request *LoginRequest) ToDTO() *dtos.LoginDTO {
	return &dtos.LoginDTO{
		Email:    request.Email,
		Password: request.Password,
	}
}
