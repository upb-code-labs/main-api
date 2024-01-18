package implementations

import (
	"context"
	"database/sql"
	"time"

	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/entities"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/errors"
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

func (repository *LaboratoriesPostgresRepository) GetLaboratoryByUUID(dto *dtos.GetLaboratoryDTO) (laboratory *entities.Laboratory, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Get base laboratory data
	query := `
		SELECT id, course_id, rubric_id, name, opening_date, due_date
		FROM laboratories
		WHERE id = $1
	`

	row := repository.Connection.QueryRowContext(ctx, query, dto.LaboratoryUUID)
	laboratory = &entities.Laboratory{}
	rubricUUID := sql.NullString{}
	if err := row.Scan(&laboratory.UUID, &laboratory.CourseUUID, &rubricUUID, &laboratory.Name, &laboratory.OpeningDate, &laboratory.DueDate); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.LaboratoryNotFoundError{}
		}

		return nil, err
	}

	if rubricUUID.Valid {
		laboratory.RubricUUID = &rubricUUID.String
	} else {
		laboratory.RubricUUID = nil
	}

	// Get markdown blocks
	markdownBlocks, err := repository.getMarkdownBlocks(dto.LaboratoryUUID)
	if err != nil {
		return nil, err
	}

	laboratory.MarkdownBlocks = markdownBlocks

	// Get test blocks
	testBlocks, err := repository.getTestBlocks(dto)
	if err != nil {
		return nil, err
	}

	laboratory.TestBlocks = testBlocks

	return laboratory, nil
}

func (repository *LaboratoriesPostgresRepository) GetLaboratoryInformationByUUID(uuid string) (laboratory *dtos.LaboratoryDetailsDTO, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Get base laboratory data
	query := `
		SELECT id, rubric_id, course_id, name, opening_date, due_date
		FROM laboratories
		WHERE id = $1
	`

	row := repository.Connection.QueryRowContext(ctx, query, uuid)
	laboratoryDetails := &dtos.LaboratoryDetailsDTO{}
	rubricUUID := sql.NullString{}

	// Parse the row
	if err := row.Scan(
		&laboratoryDetails.UUID,
		&rubricUUID,
		&laboratoryDetails.CourseUUID,
		&laboratoryDetails.Name,
		&laboratoryDetails.OpeningDate,
		&laboratoryDetails.DueDate,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.LaboratoryNotFoundError{}
		}

		return nil, err
	}

	if rubricUUID.Valid {
		laboratoryDetails.RubricUUID = &rubricUUID.String
	} else {
		laboratoryDetails.RubricUUID = nil
	}

	return laboratoryDetails, nil
}

func (repository *LaboratoriesPostgresRepository) getMarkdownBlocks(laboratoryUUID string) ([]entities.MarkdownBlock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT mb.id, mb.content, bi.block_position
		FROM markdown_blocks mb
		RIGHT JOIN blocks_index bi ON mb.block_index_id = bi.id
		WHERE mb.laboratory_id = $1
		ORDER BY bi.block_position ASC
	`

	rows, err := repository.Connection.QueryContext(ctx, query, laboratoryUUID)
	if err != nil {
		return nil, err
	}

	markdownBlocks := []entities.MarkdownBlock{}
	for rows.Next() {
		markdownBlock := entities.MarkdownBlock{}
		if err := rows.Scan(&markdownBlock.UUID, &markdownBlock.Content, &markdownBlock.Index); err != nil {
			return nil, err
		}

		markdownBlocks = append(markdownBlocks, markdownBlock)
	}

	return markdownBlocks, nil
}

func (repository *LaboratoriesPostgresRepository) getTestBlocks(dto *dtos.GetLaboratoryDTO) ([]entities.TestBlock, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT tb.id, tb.language_id, tb.test_archive_id, tb.name, bi.block_position, s.id
		FROM test_blocks tb
		RIGHT JOIN blocks_index bi ON tb.block_index_id = bi.id
		LEFT JOIN submissions s ON tb.id = s.test_block_id AND s.student_id = $2
		WHERE tb.laboratory_id = $1
		ORDER BY bi.block_position ASC
	`

	rows, err := repository.Connection.QueryContext(
		ctx,
		query,
		dto.LaboratoryUUID, dto.UserUUID,
	)
	if err != nil {
		return nil, err
	}

	testBlocks := []entities.TestBlock{}
	for rows.Next() {
		testBlock := entities.TestBlock{}
		if err := rows.Scan(
			&testBlock.UUID,
			&testBlock.LanguageUUID,
			&testBlock.TestArchiveUUID,
			&testBlock.Name,
			&testBlock.Index,
			&testBlock.SubmissionUUID,
		); err != nil {
			return nil, err
		}

		// If the user is an student, hide the test archive UUID
		if dto.UserRole == "student" {
			testBlock.TestArchiveUUID = nil
		}

		// If the user is a teacher, hide the submission UUID
		if dto.UserRole == "teacher" {
			testBlock.SubmissionUUID = nil
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

	return repository.GetLaboratoryByUUID(&dtos.GetLaboratoryDTO{
		LaboratoryUUID: laboratoryUUID,
		UserUUID:       dto.TeacherUUID,
		UserRole:       "teacher",
	})
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

func (repository *LaboratoriesPostgresRepository) CreateMarkdownBlock(laboratoryUUID string) (blockUUID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Start transaction
	tx, err := repository.Connection.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}
	defer tx.Rollback()

	// Create block index
	query := `
		INSERT INTO blocks_index (laboratory_id, block_position)
		VALUES (
			$1, 
			( SELECT COALESCE(MAX(block_position), 0) + 1 FROM blocks_index WHERE laboratory_id = $1 )
		)
		RETURNING id
	`

	row := tx.QueryRowContext(ctx, query, laboratoryUUID)
	var blockIndexUUID string
	if err := row.Scan(&blockIndexUUID); err != nil {
		return "", err
	}

	// Create markdown block
	query = `
		INSERT INTO markdown_blocks (laboratory_id, block_index_id)
		VALUES ($1, $2)
		RETURNING id
	`

	row = tx.QueryRowContext(ctx, query, laboratoryUUID, blockIndexUUID)
	if err := row.Scan(&blockUUID); err != nil {
		return "", err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return "", err
	}

	// Return the new block UUID
	return blockUUID, nil
}

func (repository *LaboratoriesPostgresRepository) CreateTestBlock(dto *dtos.CreateTestBlockDTO) (blockUUID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Start transaction
	tx, err := repository.Connection.BeginTx(ctx, nil)
	if err != nil {
		return "nil", err
	}
	defer tx.Rollback()

	// Create block index
	query := `
		INSERT INTO blocks_index (laboratory_id, block_position)
		VALUES (
			$1, 
			( SELECT COALESCE(MAX(block_position), 0) + 1 FROM blocks_index WHERE laboratory_id = $1 )
		)
		RETURNING id
	`

	row := tx.QueryRowContext(ctx, query, dto.LaboratoryUUID)
	var dbBlockIndexUUID string
	if err := row.Scan(&dbBlockIndexUUID); err != nil {
		return "", err
	}

	// Save the archive metadata
	query = `
		INSERT INTO archives (file_id)
		VALUES ($1)
		RETURNING id
	`

	row = tx.QueryRowContext(ctx, query, dto.TestArchiveUUID)
	var testBlockArchiveUUID string

	if err := row.Scan(&testBlockArchiveUUID); err != nil {
		return "", err
	}

	// Create test block
	query = `
		INSERT INTO test_blocks (language_id, test_archive_id, laboratory_id, block_index_id, name)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	row = tx.QueryRowContext(
		ctx,
		query,
		dto.LanguageUUID,
		testBlockArchiveUUID,
		dto.LaboratoryUUID,
		dbBlockIndexUUID,
		dto.Name,
	)

	var createdTestBlockUUID string
	if err := row.Scan(&createdTestBlockUUID); err != nil {
		return "", err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return "", err
	}

	// Return the new block UUID
	return createdTestBlockUUID, nil
}

func (repository *LaboratoriesPostgresRepository) GetTotalTestBlocks(laboratoryUUID string) (total int, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT COUNT(tb.id)
		FROM test_blocks AS tb
		WHERE tb.laboratory_id = $1
	`

	row := repository.Connection.QueryRowContext(ctx, query, laboratoryUUID)
	if err := row.Scan(&total); err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.LaboratoryNotFoundError{}
		}

		return 0, err
	}

	return total, nil
}

func (repository *LaboratoriesPostgresRepository) GetStudentsProgress(laboratoryUUID string) (progress []*dtos.SummarizedStudentProgressDTO, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT spv.student_id, spv.student_full_name, spv.pending_submissions, spv.running_submissions, spv.failing_submissions, spv.success_submissions
		FROM students_progress_view AS spv
		WHERE spv.laboratory_id = $1
	`

	rows, err := repository.Connection.QueryContext(ctx, query, laboratoryUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.LaboratoryNotFoundError{}
		}

		return nil, err
	}

	progress = []*dtos.SummarizedStudentProgressDTO{}

	for rows.Next() {
		studentProgress := dtos.SummarizedStudentProgressDTO{}

		if err := rows.Scan(
			&studentProgress.StudentUUID,
			&studentProgress.StudentFullName,
			&studentProgress.PendingSubmissions,
			&studentProgress.RunningSubmissions,
			&studentProgress.FailingSubmissions,
			&studentProgress.SuccessSubmissions,
		); err != nil {
			return nil, err
		}

		progress = append(progress, &studentProgress)
	}

	return progress, nil
}

// DoesTeacherOwnLaboratory returns true if the teacher owns the laboratory
// and throws an error if the laboratory does not exist
func (repository *LaboratoriesPostgresRepository) DoesTeacherOwnLaboratory(teacherUUID string, laboratoryUUID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT l.id, c.teacher_id
		FROM laboratories AS l
		INNER JOIN courses AS c ON l.course_id = c.id
		WHERE l.id = $1
	`

	row := repository.Connection.QueryRowContext(ctx, query, laboratoryUUID)

	var laboratoryID, teacherID string
	if err := row.Scan(&laboratoryID, &teacherID); err != nil {
		if err == sql.ErrNoRows {
			return false, errors.LaboratoryNotFoundError{}
		}

		return false, err
	}

	return teacherID == teacherUUID, nil
}

// GetStudentSubmissions returns the submissions of a student in a laboratory
func (repository *LaboratoriesPostgresRepository) GetStudentSubmissions(laboratoryUUID string, studentUUID string) (submissions []*dtos.SummarizedStudentSubmissionDTO, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT s.id, s.archive_id, tb.name, s.status, s.passing
		FROM submissions AS s
		INNER JOIN test_blocks AS tb ON s.test_block_id = tb.id
		WHERE tb.laboratory_id = $1 AND s.student_id = $2
	`

	rows, err := repository.Connection.QueryContext(ctx, query, laboratoryUUID, studentUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return an empty state if the student has not submitted anything
			return []*dtos.SummarizedStudentSubmissionDTO{}, nil
		}

		return nil, err
	}

	submissions = []*dtos.SummarizedStudentSubmissionDTO{}
	for rows.Next() {
		submission := dtos.SummarizedStudentSubmissionDTO{}

		if err := rows.Scan(
			&submission.SubmissionUUID,
			&submission.SubmissionArchiveUUID,
			&submission.TestBlockName,
			&submission.SubmissionStatus,
			&submission.IsSubmissionPassing,
		); err != nil {
			return nil, err
		}

		submissions = append(submissions, &submission)
	}

	return submissions, nil
}
