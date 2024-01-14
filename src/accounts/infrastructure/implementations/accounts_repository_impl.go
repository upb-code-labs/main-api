package implementations

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
		INSERT INTO users (role, email, full_name, password_hash, created_by)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := repository.Connection.ExecContext(
		ctx,
		query,
		"admin",
		dto.Email,
		dto.FullName,
		dto.Password,
		dto.CreatedBy,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repository *AccountsPostgresRepository) SaveTeacher(dto dtos.RegisterUserDTO) error {
	// Save user
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (role, email, full_name, password_hash, created_by)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := repository.Connection.ExecContext(
		ctx,
		query,
		"teacher",
		dto.Email,
		dto.FullName,
		dto.Password,
		dto.CreatedBy,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repository *AccountsPostgresRepository) GetUserByUUID(uuid string) (*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT id, role, institutional_id, email, full_name, password_hash
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	row := repository.Connection.QueryRowContext(ctx, query, uuid)
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

func (repository *AccountsPostgresRepository) GetAdmins() ([]*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT id, institutional_id, email, full_name, created_at, creator_full_name
		FROM users_with_creator
		WHERE role = 'admin'
	`

	rows, err := repository.Connection.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var admins []*entities.User
	for rows.Next() {
		var admin entities.User
		var adminInstitutionalId sql.NullString
		err := rows.Scan(
			&admin.UUID,
			&adminInstitutionalId,
			&admin.Email,
			&admin.FullName,
			&admin.CreatedAt,
			&admin.CreatedBy,
		)

		if err != nil {
			return nil, err
		}

		if adminInstitutionalId.Valid {
			admin.InstitutionalId = adminInstitutionalId.String
		}

		admins = append(admins, &admin)
	}

	return admins, nil
}

func (repository *AccountsPostgresRepository) SearchStudentsByFullName(fullName string) ([]*entities.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT id, institutional_id, email, full_name
		FROM users
		WHERE role = 'student' AND lower(full_name) LIKE lower($1)
	`

	rows, err := repository.Connection.QueryContext(ctx, query, fullName+"%")
	if err != nil {
		return nil, err
	}

	var students []*entities.User
	for rows.Next() {
		var student entities.User
		var studentInstitutionalId sql.NullString
		err := rows.Scan(
			&student.UUID,
			&studentInstitutionalId,
			&student.Email,
			&student.FullName,
		)

		if err != nil {
			return nil, err
		}

		if studentInstitutionalId.Valid {
			student.InstitutionalId = studentInstitutionalId.String
		}

		students = append(students, &student)
	}

	return students, nil
}

func (repository *AccountsPostgresRepository) UpdatePassword(uuid string, newPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		UPDATE users
		SET password_hash = $1
		WHERE id = $2
	`

	_, err := repository.Connection.ExecContext(ctx, query, newPassword, uuid)
	if err != nil {
		return err
	}

	return nil
}

// UpdateProfile updates the profile of a user in the database
func (repository *AccountsPostgresRepository) UpdateProfile(dto dtos.UpdateAccountDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		UPDATE users
		SET full_name = $1, email = $2, institutional_id = $3
		WHERE id = $4
	`

	_, err := repository.Connection.ExecContext(
		ctx,
		query,
		dto.FullName,
		dto.Email,
		dto.InstitutionalId,
		dto.UserUUID,
	)
	if err != nil {
		return err
	}

	return nil
}
