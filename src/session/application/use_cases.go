package application

import (
	"database/sql"

	accountsDefinitions "github.com/UPB-Code-Labs/main-api/src/accounts/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/accounts/domain/entities"
	accountsErrors "github.com/UPB-Code-Labs/main-api/src/accounts/domain/errors"
	"github.com/UPB-Code-Labs/main-api/src/session/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/session/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/session/domain/errors"
)

type SessionUseCases struct {
	AccountsRepository accountsDefinitions.AccountsRepository
	PasswordHasher     accountsDefinitions.PasswordsHasher
	TokenHandler       definitions.TokenHandler
}

func (useCases *SessionUseCases) Login(dto dtos.LoginDTO) (dtos.SessionDTO, error) {
	sessionResponse := dtos.SessionDTO{}

	// Get the user
	user, err := useCases.AccountsRepository.GetUserByEmail(dto.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return sessionResponse, errors.InvalidCredentialsError{}
		}

		return sessionResponse, err
	}

	// Check the password
	valid, err := useCases.PasswordHasher.ComparePasswords(dto.Password, user.PasswordHash)
	if err != nil {
		return sessionResponse, err
	}
	if !valid {
		return sessionResponse, errors.InvalidCredentialsError{}
	}

	// Generate the token
	token, err := useCases.TokenHandler.GenerateToken(*user)
	if err != nil {
		return sessionResponse, err
	}

	// Return the session
	sessionResponse.User = *user
	sessionResponse.Token = token
	return sessionResponse, nil
}

func (useCases *SessionUseCases) WhoAmI(uuid string) (entities.User, error) {
	user, err := useCases.AccountsRepository.GetUserByUUID(uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.User{}, accountsErrors.UserNotFoundError{
				Uuuid: uuid,
			}
		}

		return entities.User{}, err
	}

	return *user, nil
}
