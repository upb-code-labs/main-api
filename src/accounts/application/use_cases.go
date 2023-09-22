package application

import (
	"database/sql"

	"github.com/UPB-Code-Labs/main-api/src/accounts/domain"
	"github.com/UPB-Code-Labs/main-api/src/accounts/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/accounts/domain/errors"
)

type AccountsUseCases struct {
	AccountsRepository domain.AccountsRepository
	PasswordsHasher    domain.PasswordsHasher
}

func (useCases *AccountsUseCases) RegisterStudent(dto dtos.RegisterUserDTO) error {
	// Check if email is already in use
	existingUser, err := useCases.AccountsRepository.GetUserByEmail(dto.Email)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if existingUser != nil {
		return errors.EmailAlreadyInUseError{Email: dto.Email}
	}

	// Check if institutional ID is already in use
	existingUser, err = useCases.AccountsRepository.GetUserByInstitutionalId(dto.InstitutionalId)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if existingUser != nil {
		return errors.InstitutionalIdAlreadyInUseError{InstitutionalId: dto.InstitutionalId}
	}

	// Hash password
	hash, err := useCases.PasswordsHasher.HashPassword(dto.Password)
	if err != nil {
		return err
	}
	dto.Password = hash

	// Save user
	err = useCases.AccountsRepository.SaveStudent(dto)
	return err
}

func (useCases *AccountsUseCases) RegisterAdmin(dto dtos.RegisterUserDTO) error {
	// Check if email is already in use
	existingUser, err := useCases.AccountsRepository.GetUserByEmail(dto.Email)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if existingUser != nil {
		return errors.EmailAlreadyInUseError{Email: dto.Email}
	}

	// Hash password
	hash, err := useCases.PasswordsHasher.HashPassword(dto.Password)
	if err != nil {
		return err
	}
	dto.Password = hash

	// Save user
	err = useCases.AccountsRepository.SaveAdmin(dto)
	return err
}
