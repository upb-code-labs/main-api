package infrastructure

import (
	"context"
	"database/sql"
	"time"

	"github.com/UPB-Code-Labs/main-api/src/accounts/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/accounts/domain/entities"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
)

type AccountsPostgresRepository struct {
	Connection *sql.DB
}

var accountsRepositoryInstance *AccountsPostgresRepository

func GetAccountsPgRepository() *AccountsPostgresRepository {
	if accountsRepositoryInstance == nil {
		accountsRepositoryInstance = &AccountsPostgresRepository{
			Connection: infrastructure.GetPostgresConnection(),
		}
	}

	return accountsRepositoryInstance
}

func (repository *AccountsPostgresRepository) SaveStudent(dto dtos.RegisterUserDTO) error {
	// Save user
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (role, institutional_id, email, full_name, password_hash)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := repository.Connection.ExecContext(
		ctx,
		query,
		"student",
		dto.InstitutionalId,
		dto.Email,
		dto.FullName,
		dto.Password,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repository *AccountsPostgresRepository) SaveAdmin(dto dtos.RegisterUserDTO) error {
	// Save user
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (role, email, full_name, password_hash)
		VALUES ($1, $2, $3, $4)
	`

	_, err := repository.Connection.ExecContext(
		ctx,
		query,
		"admin",
		dto.Email,
		dto.FullName,
		dto.Password,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repository *AccountsPostgresRepository) GetUserByEmail(email string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT id, role, institutional_id, email, full_name, password_hash
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	row := repository.Connection.QueryRowContext(ctx, query, email)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var user entities.User
	var userInstitutionalId sql.NullString
	err := row.Scan(
		&user.UUID,
		&user.Role,
		&userInstitutionalId,
		&user.Email,
		&user.FullName,
		&user.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	if userInstitutionalId.Valid {
		user.InstitutionalId = userInstitutionalId.String
	}

	return &user, nil
}

func (repository *AccountsPostgresRepository) GetUserByInstitutionalId(id string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT id, role, institutional_id, email, full_name, password_hash
		FROM users
		WHERE institutional_id = $1
		LIMIT 1
	`

	row := repository.Connection.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var user entities.User
	var userInstitutionalId sql.NullString
	err := row.Scan(
		&user.UUID,
		&user.Role,
		&userInstitutionalId,
		&user.Email,
		&user.FullName,
		&user.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	if userInstitutionalId.Valid {
		user.InstitutionalId = userInstitutionalId.String
	}

	return &user, nil
}
