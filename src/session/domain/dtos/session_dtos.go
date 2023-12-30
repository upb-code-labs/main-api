package dtos

import "github.com/UPB-Code-Labs/main-api/src/accounts/domain/entities"

type CustomClaimsDTO struct {
	UUID string
	Role string
}

type SessionDTO struct {
	User  entities.User
	Token string
}

type LoginDTO struct {
	Email    string
	Password string
}
