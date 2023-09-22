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

var rolesCache map[string]string = make(map[string]string)

func (repository *AccountsPostgresRepository) getRoleUUID(role string) (string, error) {
	if _, ok := rolesCache[role]; !ok {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		query := `
			SELECT id
			FROM roles
			WHERE name = $1
			LIMIT 1
		`

		row := repository.Connection.QueryRowContext(ctx, query, role)
		if row.Err() != nil {
			return "", row.Err()
		}

		var uuid string
		err := row.Scan(&uuid)
		if err != nil {
			return "", err
		}

		rolesCache[role] = uuid
	}

	return rolesCache[role], nil
}

func (repository *AccountsPostgresRepository) SaveStudent(dto dtos.RegisterUserDTO) error {
	// Get role UUID
	studentRole, err := repository.getRoleUUID("student")
	if err != nil {
		return err
	}

	// Save user
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (role_id, institutional_id, email, full_name, password_hash)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = repository.Connection.ExecContext(
		ctx,
		query,
		studentRole,
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

func (repository *AccountsPostgresRepository) GetUserByEmail(email string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT id, role_id, institutional_id, email, full_name, password_hash
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	row := repository.Connection.QueryRowContext(ctx, query, email)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var user entities.User
	err := row.Scan(
		&user.UUID,
		&user.RoleUUID,
		&user.InstitutionalId,
		&user.Email,
		&user.FullName,
		&user.PasswordHash,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repository *AccountsPostgresRepository) GetUserByInstitutionalId(id string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT id, role_id, institutional_id, email, full_name, password_hash
		FROM users
		WHERE institutional_id = $1
		LIMIT 1
	`

	row := repository.Connection.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var user entities.User
	err := row.Scan(
		&user.UUID,
		&user.RoleUUID,
		&user.InstitutionalId,
		&user.Email,
		&user.FullName,
		&user.PasswordHash,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
