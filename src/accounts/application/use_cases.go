package application

import (
	"database/sql"

	"github.com/UPB-Code-Labs/main-api/src/accounts/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/accounts/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/accounts/domain/entities"
	"github.com/UPB-Code-Labs/main-api/src/accounts/domain/errors"
	sessionErrors "github.com/UPB-Code-Labs/main-api/src/session/domain/errors"
)

type AccountsUseCases struct {
	AccountsRepository definitions.AccountsRepository
	PasswordsHasher    definitions.PasswordsHasher
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

func (useCases *AccountsUseCases) RegisterTeacher(dto dtos.RegisterUserDTO) error {
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
	err = useCases.AccountsRepository.SaveTeacher(dto)
	return err
}

func (useCases *AccountsUseCases) GetAdmins() ([]*entities.User, error) {
	return useCases.AccountsRepository.GetAdmins()
}

func (useCases *AccountsUseCases) SearchStudentsByFullName(fullName string) ([]*entities.User, error) {
	return useCases.AccountsRepository.SearchStudentsByFullName(fullName)
}

func (useCases *AccountsUseCases) UpdatePassword(dto dtos.UpdatePasswordDTO) error {
	// Get user
	user, err := useCases.AccountsRepository.GetUserByUUID(dto.UserUUID)
	if err != nil {
		return err
	}

	// Check if old password is correct
	doesOldPasswordMatch, err := useCases.PasswordsHasher.ComparePasswords(dto.OldPassword, user.PasswordHash)
	if err != nil {
		return err
	}

	if !doesOldPasswordMatch {
		return sessionErrors.InvalidCredentialsError{}
	}

	// Hash new password
	newPasswordHash, err := useCases.PasswordsHasher.HashPassword(dto.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	err = useCases.AccountsRepository.UpdatePassword(dto.UserUUID, newPasswordHash)
	if err != nil {
		return err
	}

	return nil
}

func (useCases *AccountsUseCases) UpdateProfile(dto dtos.UpdateAccountDTO) error {
	// Get user
	user, err := useCases.AccountsRepository.GetUserByUUID(dto.UserUUID)
	if err != nil {
		return err
	}

	// Check if the given password is correct
	doesPasswordMatch, err := useCases.PasswordsHasher.ComparePasswords(dto.Password, user.PasswordHash)
	if err != nil {
		return err
	}

	if !doesPasswordMatch {
		return sessionErrors.InvalidCredentialsError{}
	}

	// Check if email is already in use
	existingUser, err := useCases.AccountsRepository.GetUserByEmail(dto.Email)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if existingUser != nil {
		return errors.EmailAlreadyInUseError{Email: dto.Email}
	}

	// Check if institutional ID is already in use
	existingUser, err = useCases.AccountsRepository.GetUserByInstitutionalId(*dto.InstitutionalId)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if existingUser != nil {
		return errors.InstitutionalIdAlreadyInUseError{InstitutionalId: *dto.InstitutionalId}
	}

	// Update profile
	err = useCases.AccountsRepository.UpdateProfile(dto)
	return err
}
