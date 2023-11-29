package implementations

import (
	"context"
	"database/sql"
	"time"

	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
)

type BlocksPostgresRepository struct {
	Connection *sql.DB
}

// Singleton
var blocksPostgresRepositoryInstance *BlocksPostgresRepository

func GetBlocksPostgresRepositoryInstance() *BlocksPostgresRepository {
	if blocksPostgresRepositoryInstance == nil {
		blocksPostgresRepositoryInstance = &BlocksPostgresRepository{
			Connection: infrastructure.GetPostgresConnection(),
		}
	}

	return blocksPostgresRepositoryInstance
}

func (repository *BlocksPostgresRepository) UpdateMarkdownBlockContent(blockUUID string, content string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		UPDATE markdown_blocks
		SET content = $1
		WHERE id = $2
	`

	_, err = repository.Connection.ExecContext(ctx, query, content, blockUUID)
	return err
}

func (repository *BlocksPostgresRepository) DoesTeacherOwnsMarkdownBlock(teacherUUID string, blockUUID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Get the UUID of the laboratory the block belongs to
	query := `
		SELECT laboratory_id
		FROM markdown_blocks
		WHERE id = $1
	`

	row := repository.Connection.QueryRowContext(ctx, query, blockUUID)
	var laboratoryUUID string
	if err := row.Scan(&laboratoryUUID); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
	}

	// Check if the teacher owns the laboratory
	query = `
		SELECT teacher_id
		FROM courses
		WHERE id = (
			SELECT course_id
			FROM laboratories
			WHERE id = $1
		)
	`

	row = repository.Connection.QueryRowContext(ctx, query, laboratoryUUID)
	var laboratoryTeacherUUID string
	if err := row.Scan(&laboratoryTeacherUUID); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
	}

	return laboratoryTeacherUUID == teacherUUID, nil
}
