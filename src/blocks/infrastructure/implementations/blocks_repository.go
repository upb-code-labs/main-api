package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/errors"
	laboratoriesEntities "github.com/UPB-Code-Labs/main-api/src/laboratories/domain/entities"
	laboratoriesDomainErrors "github.com/UPB-Code-Labs/main-api/src/laboratories/domain/errors"
	sharedEntities "github.com/UPB-Code-Labs/main-api/src/shared/domain/entities"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	staticFilesDefinitions "github.com/UPB-Code-Labs/main-api/src/static-files/domain/definitions"
	staticFilesDTOs "github.com/UPB-Code-Labs/main-api/src/static-files/domain/dtos"
	staticFilesImplementations "github.com/UPB-Code-Labs/main-api/src/static-files/infrastructure/implementations"
)

type BlocksPostgresRepository struct {
	Connection            *sql.DB
	StaticFilesRepository staticFilesDefinitions.StaticFilesRepository
}

// Singleton
var blocksPostgresRepositoryInstance *BlocksPostgresRepository

func GetBlocksPostgresRepositoryInstance() *BlocksPostgresRepository {
	if blocksPostgresRepositoryInstance == nil {
		blocksPostgresRepositoryInstance = &BlocksPostgresRepository{
			Connection:            sharedInfrastructure.GetPostgresConnection(),
			StaticFilesRepository: &staticFilesImplementations.StaticFilesMicroserviceImplementation{},
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
			return false, errors.BlockNotFound{}
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
			return false, laboratoriesDomainErrors.LaboratoryNotFoundError{}
		}
	}

	return laboratoryTeacherUUID == teacherUUID, nil
}

func (repository *BlocksPostgresRepository) DoesTeacherOwnsTestBlock(teacherUUID string, blockUUID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Get the UUID of the laboratory the block belongs to
	query := `
		SELECT laboratory_id
		FROM test_blocks
		WHERE id = $1
	`

	row := repository.Connection.QueryRowContext(ctx, query, blockUUID)
	var laboratoryUUID string
	if err := row.Scan(&laboratoryUUID); err != nil {
		if err == sql.ErrNoRows {
			return false, &errors.BlockNotFound{}
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
			return false, laboratoriesDomainErrors.LaboratoryNotFoundError{}
		}
	}

	return laboratoryTeacherUUID == teacherUUID, nil
}

func (repository *BlocksPostgresRepository) CanStudentSubmitToTestBlock(studentUUID string, testBlockUUID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Check if the student is enrolled in the laboratory the block belongs to
	query := `
		SELECT user_id
		FROM courses_has_users
		WHERE course_id = (
			SELECT id
			FROM courses
			WHERE id = (
				SELECT course_id
				FROM laboratories
				WHERE id = (
					SELECT laboratory_id
					FROM test_blocks
					WHERE id = $1
				)
			)
		) 
		AND user_id = $2 
		AND is_user_active = true
	`

	row := repository.Connection.QueryRowContext(ctx, query, testBlockUUID, studentUUID)
	var studentID string
	if err := row.Scan(&studentID); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (repository *BlocksPostgresRepository) GetTestArchiveUUIDFromTestBlockUUID(blockUUID string) (uuid string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT file_id
		FROM archives
		WHERE id = (
			SELECT test_archive_id
			FROM test_blocks
			WHERE id = $1
		)
	`

	row := repository.Connection.QueryRowContext(ctx, query, blockUUID)

	// Parse the row
	err = row.Scan(&uuid)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", &errors.BlockNotFound{}
		}

		return "", err
	}

	return uuid, nil
}

func (repository *BlocksPostgresRepository) UpdateTestBlock(dto *dtos.UpdateTestBlockDTO) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Update the block
	query := `
		UPDATE test_blocks
		SET language_id = $1, name = $2
		WHERE id = $3
	`

	_, err = repository.Connection.ExecContext(ctx, query, dto.LanguageUUID, dto.Name, dto.BlockUUID)
	if err != nil {
		return err
	}

	return nil
}

func (repository *BlocksPostgresRepository) GetTestBlockLaboratoryUUID(blockUUID string) (laboratoryUUID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT laboratory_id
		FROM test_blocks
		WHERE id = $1
	`

	row := repository.Connection.QueryRowContext(ctx, query, blockUUID)
	if err := row.Scan(&laboratoryUUID); err != nil {
		return "", err
	}

	return laboratoryUUID, nil
}

func (repository *BlocksPostgresRepository) DeleteMarkdownBlock(blockUUID string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Get the UUID of the block index
	query := `
		SELECT block_index_id
		FROM markdown_blocks
		WHERE id = $1
	`

	row := repository.Connection.QueryRowContext(ctx, query, blockUUID)
	var blockIndexUUID string
	if err := row.Scan(&blockIndexUUID); err != nil {
		if err == sql.ErrNoRows {
			return &errors.BlockNotFound{}
		}

		return err
	}

	// After deleting the block index, the block will be deleted automatically due to the `ON DELETE CASCADE` constraint
	err = repository.deleteBlockIndex(blockIndexUUID)
	if err != nil {
		return err
	}

	return nil
}

func (repository *BlocksPostgresRepository) DeleteTestBlock(blockUUID string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Get the UUIDs of the dependent archives before deleting the block
	dependentArchivesUUIDs, err := repository.getDependentArchivesByTestBlockUUID(blockUUID)
	if err != nil {
		return err
	}

	// Delete the dependent archives in a separate goroutine
	go repository.deleteDependentArchives(dependentArchivesUUIDs)

	// Get the UUID of the block index
	query := `
		SELECT block_index_id
		FROM test_blocks
		WHERE id = $1
	`

	row := repository.Connection.QueryRowContext(ctx, query, blockUUID)
	var blockIndexUUID string
	if err := row.Scan(&blockIndexUUID); err != nil {
		if err == sql.ErrNoRows {
			return &errors.BlockNotFound{}
		}

		return err
	}

	// After deleting the block index, the block will be deleted automatically due to the `ON DELETE CASCADE` constraint
	err = repository.deleteBlockIndex(blockIndexUUID)
	if err != nil {
		return err
	}

	return nil
}

func (repository *BlocksPostgresRepository) deleteBlockIndex(blockIndexUUID string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		DELETE FROM blocks_index
		WHERE id = $1
	`

	_, err = repository.Connection.ExecContext(ctx, query, blockIndexUUID)
	if err != nil {
		return err
	}

	return nil
}

func (repository *BlocksPostgresRepository) getDependentArchivesByTestBlockUUID(blockUUID string) (archives []*sharedEntities.StaticFileArchive, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Get the UUID of the test block's tests archive
	query := `
		SELECT file_id
		FROM archives
		WHERE id = (
			SELECT test_archive_id
			FROM test_blocks
			WHERE id = $1
		)
	`

	row := repository.Connection.QueryRowContext(ctx, query, blockUUID)
	var testsArchiveUUID string
	if err := row.Scan(&testsArchiveUUID); err != nil {
		if err == sql.ErrNoRows {
			return nil, &errors.BlockNotFound{}
		}

		return nil, err
	}

	archives = append(archives, &sharedEntities.StaticFileArchive{
		ArchiveUUID: testsArchiveUUID,
		ArchiveType: "test",
	})

	// Get the UUID of the test block's submissions archives
	query = `
		SELECT file_id
		FROM archives
		WHERE id IN (
			SELECT archive_id
			FROM submissions
			WHERE test_block_id = $1
		)
	`

	rows, err := repository.Connection.QueryContext(ctx, query, blockUUID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var submissionArchiveUUID string
		if err := rows.Scan(&submissionArchiveUUID); err != nil {
			return nil, err
		}

		archives = append(archives, &sharedEntities.StaticFileArchive{
			ArchiveUUID: submissionArchiveUUID,
			ArchiveType: "submission",
		})
	}

	return archives, nil
}

func (repository *BlocksPostgresRepository) deleteDependentArchives(archives []*sharedEntities.StaticFileArchive) {
	log.Printf("[INFO] - [BlocksPostgresRepository] - [deleteDependentArchives]: Deleting %d archives \n", len(archives))

	for _, archive := range archives {
		err := repository.StaticFilesRepository.DeleteArchive(
			&staticFilesDTOs.StaticFileArchiveDTO{
				FileUUID: archive.ArchiveUUID,
				FileType: archive.ArchiveType,
			},
		)
		if err != nil {
			log.Printf("[ERROR] - [BlocksPostgresRepository] - [deleteDependentArchives]: %s \n", err.Error())
		}
	}
}

func (repository *BlocksPostgresRepository) GetMarkdownBlockByUUID(blockUUID string) (markdownBlock *laboratoriesEntities.MarkdownBlock, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
	SELECT mb.id, mb.content, bi.block_position
	FROM markdown_blocks mb
	RIGHT JOIN blocks_index bi ON mb.block_index_id = bi.id
	WHERE mb.id = $1
	`

	row := repository.Connection.QueryRowContext(ctx, query, blockUUID)

	// Parse the row
	markdownBlock = &laboratoriesEntities.MarkdownBlock{}
	err = row.Scan(
		&markdownBlock.UUID,
		&markdownBlock.Content,
		&markdownBlock.Index,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &errors.BlockNotFound{}
		}

		return nil, err
	}

	return markdownBlock, nil
}

func (repository *BlocksPostgresRepository) GetTestBlockByUUID(blockUUID string) (testBlock *laboratoriesEntities.TestBlock, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT tb.id, tb.language_id, tb.test_archive_id, tb.name, bi.block_position
		FROM test_blocks tb
		RIGHT JOIN blocks_index bi ON tb.block_index_id = bi.id
		WHERE tb.id = $1
	`

	row := repository.Connection.QueryRowContext(ctx, query, blockUUID)

	// Parse the row
	testBlock = &laboratoriesEntities.TestBlock{}
	err = row.Scan(
		&testBlock.UUID,
		&testBlock.LanguageUUID,
		&testBlock.TestArchiveUUID,
		&testBlock.Name,
		&testBlock.Index,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &errors.BlockNotFound{}
		}

		return nil, err
	}

	return testBlock, nil
}

func (repository *BlocksPostgresRepository) SwapBlocks(firstBlockUUID, secondBlockUUID string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Define query to call the stored procedure
	query := `
		SELECT swap_blocks_index($1, $2)
	`

	// Start a transaction
	tx, err := repository.Connection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Run the query
	fmt.Println("Swapping blocks", firstBlockUUID, secondBlockUUID)
	_, err = tx.ExecContext(ctx, query, firstBlockUUID, secondBlockUUID)
	if err != nil {
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
