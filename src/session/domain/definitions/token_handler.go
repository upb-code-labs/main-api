package definitions

import (
	"github.com/UPB-Code-Labs/main-api/src/accounts/domain/entities"
	"github.com/UPB-Code-Labs/main-api/src/session/domain/dtos"
)

type TokenHandler interface {
	GenerateToken(user entities.User) (string, error)
	ValidateToken(token string) (dtos.CustomClaimsDTO, error)
}
