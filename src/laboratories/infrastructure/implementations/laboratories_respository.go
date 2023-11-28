package implementations

import (
	"context"
	"database/sql"
	"time"

	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/entities"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
)

type LaboratoriesPostgresRepository struct {
	Connection *sql.DB
}

// Singleton
var laboratoriesPostgresRepositoryInstance *LaboratoriesPostgresRepository

func GetLaboratoriesPostgresRepositoryInstance() *LaboratoriesPostgresRepository {
	if laboratoriesPostgresRepositoryInstance == nil {
		laboratoriesPostgresRepositoryInstance = &LaboratoriesPostgresRepository{
			Connection: infrastructure.GetPostgresConnection(),
		}
	}

	return laboratoriesPostgresRepositoryInstance
}

func (repository *LaboratoriesPostgresRepository) GetLaboratoryByUUID(uuid string) (laboratory *entities.Laboratory, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Get base laboratory data
	query := `
		SELECT id, course_id, rubric_id, name, opening_date, due_date
		FROM laboratories
		WHERE id = $1
	`

	row := repository.Connection.QueryRowContext(ctx, query, uuid)
	laboratory = &entities.Laboratory{}
	rubricUUID := sql.NullString{}
	if err := row.Scan(&laboratory.UUID, &laboratory.CourseUUID, &rubricUUID, &laboratory.Name, &laboratory.OpeningDate, &laboratory.DueDate); err != nil {
		return nil, err
	}
	if rubricUUID.Valid {
		laboratory.RubricUUID = rubricUUID.String
	}

	// Get markdown blocks
	markdownBlocks, err := repository.getMarkdownBlocks(uuid)
	if err != nil {
		return nil, err
	}

	laboratory.MarkdownBlocks = markdownBlocks

	// Get test blocks
	testBlocks, err := repository.getTestBlocks(uuid)
	if err != nil {
		return nil, err
	}

	laboratory.TestBlocks = testBlocks

	return laboratory, nil
}

func (repository *LaboratoriesPostgresRepository) getMarkdownBlocks(laboratoryUUID string) ([]entities.MarkdownBlock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT id, content, block_index
		FROM markdown_blocks
		WHERE laboratory_id = $1
		ORDER BY block_index ASC
	`

	rows, err := repository.Connection.QueryContext(ctx, query, laboratoryUUID)
	if err != nil {
		return nil, err
	}

	markdownBlocks := []entities.MarkdownBlock{}
	for rows.Next() {
		markdownBlock := entities.MarkdownBlock{}
		if err := rows.Scan(&markdownBlock.UUID, &markdownBlock.Content, &markdownBlock.Order); err != nil {
			return nil, err
		}

		markdownBlocks = append(markdownBlocks, markdownBlock)
	}

	return markdownBlocks, nil
}

func (repository *LaboratoriesPostgresRepository) getTestBlocks(laboratoryUUID string) ([]entities.TestBlock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT id, language_id, tests_archive_id, name, block_index
		FROM test_blocks
		WHERE laboratory_id = $1
		ORDER BY block_index ASC
	`

	rows, err := repository.Connection.QueryContext(ctx, query, laboratoryUUID)
	if err != nil {
		return nil, err
	}

	testBlocks := []entities.TestBlock{}
	for rows.Next() {
		testBlock := entities.TestBlock{}
		if err := rows.Scan(&testBlock.UUID, &testBlock.LanguageUUID, &testBlock.TestArchiveUUID, &testBlock.Name, &testBlock.Order); err != nil {
			return nil, err
		}

		testBlocks = append(testBlocks, testBlock)
	}

	return testBlocks, nil
}

func (repository *LaboratoriesPostgresRepository) SaveLaboratory(dto *dtos.CreateLaboratoryDTO) (laboratory *entities.Laboratory, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		INSERT INTO laboratories (course_id, name, opening_date, due_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	row := repository.Connection.QueryRowContext(ctx, query, dto.CourseUUID, dto.Name, dto.OpeningDate, dto.DueDate)
	var laboratoryUUID string
	if err := row.Scan(&laboratoryUUID); err != nil {
		return nil, err
	}

	return repository.GetLaboratoryByUUID(laboratoryUUID)
}

func (repository *LaboratoriesPostgresRepository) UpdateLaboratory(dto *dtos.UpdateLaboratoryDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		UPDATE laboratories
		SET name = $1, opening_date = $2, due_date = $3, rubric_id = $4
		WHERE id = $5
	`

	_, err := repository.Connection.ExecContext(ctx, query, dto.Name, dto.OpeningDate, dto.DueDate, dto.RubricUUID, dto.LaboratoryUUID)
	return err
}
