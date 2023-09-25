package domain

import (
	"github.com/UPB-Code-Labs/main-api/src/accounts/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/accounts/domain/entities"
)

type AccountsRepository interface {
	SaveStudent(dto dtos.RegisterUserDTO) error
	SaveAdmin(dto dtos.RegisterUserDTO) error
	SaveTeacher(dto dtos.RegisterUserDTO) error

	GetUserByEmail(email string) (*entities.User, error)
	GetUserByInstitutionalId(institutionalId string) (*entities.User, error)
}
