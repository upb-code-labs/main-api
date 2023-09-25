package domain

import "github.com/UPB-Code-Labs/main-api/src/accounts/domain/entities"

type JwtCustomClaims struct {
	UUID string
	Role string
}

type TokenHandler interface {
	GenerateToken(user entities.User) (string, error)
	ValidateToken(token string) (JwtCustomClaims, error)
}
