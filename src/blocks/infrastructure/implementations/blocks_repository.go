package implementations

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/errors"
	laboratoriesDomainErrors "github.com/UPB-Code-Labs/main-api/src/laboratories/domain/errors"
	sharedEntities "github.com/UPB-Code-Labs/main-api/src/shared/domain/entities"
	sharedDomainErrors "github.com/UPB-Code-Labs/main-api/src/shared/domain/errors"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
)

type BlocksPostgresRepository struct {
	Connection *sql.DB
}

// Singleton
var blocksPostgresRepositoryInstance *BlocksPostgresRepository

func GetBlocksPostgresRepositoryInstance() *BlocksPostgresRepository {
	if blocksPostgresRepositoryInstance == nil {
		blocksPostgresRepositoryInstance = &BlocksPostgresRepository{
			Connection: sharedInfrastructure.GetPostgresConnection(),
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

func (repository *BlocksPostgresRepository) SaveTestsArchive(file *multipart.File) (uuid string, err error) {
	// Create multipart writer
	staticFilesEndpoint := fmt.Sprintf("%s/archives/save", sharedInfrastructure.GetEnvironment().StaticFilesMicroserviceAddress)
	baseMultipartBuffer, err := sharedInfrastructure.GetMultipartFormBufferFromFile(file)
	if err != nil {
		return "", err
	}

	// Add the file type field to the request
	err = baseMultipartBuffer.BodyBufferWriter.WriteField("archive_type", "test")
	if err != nil {
		return "", err
	}

	// Close the writer
	err = baseMultipartBuffer.BodyBufferWriter.Close()
	if err != nil {
		return "", err
	}

	// Prepare the request
	req, err := http.NewRequest("POST", staticFilesEndpoint, baseMultipartBuffer.BodyBuffer)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", baseMultipartBuffer.BodyBufferWriter.FormDataContentType())

	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)

	// Forward error message if any
	microserviceError := sharedInfrastructure.ParseMicroserviceError(res, err)
	if microserviceError != nil {
		return "", microserviceError
	}

	// Return the UUID of the saved file
	var response map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	if response["uuid"] == nil {
		return "", &sharedDomainErrors.GenericDomainError{
			Code:    http.StatusInternalServerError,
			Message: "The static files microservice did not return the UUID of the saved file",
		}
	}

	return response["uuid"].(string), nil
}

func (repository *BlocksPostgresRepository) OverwriteTestsArchive(uuid string, file *multipart.File) (err error) {
	// Create multipart writer
	staticFilesEndpoint := fmt.Sprintf("%s/archives/overwrite", sharedInfrastructure.GetEnvironment().StaticFilesMicroserviceAddress)
	baseMultipartBuffer, err := sharedInfrastructure.GetMultipartFormBufferFromFile(file)
	if err != nil {
		return err
	}

	// Add the file type field to the request
	err = baseMultipartBuffer.BodyBufferWriter.WriteField("archive_type", "test")
	if err != nil {
		return err
	}

	// Add the archive uuid field to the request
	err = baseMultipartBuffer.BodyBufferWriter.WriteField("archive_uuid", uuid)
	if err != nil {
		return err
	}

	// Close the writer
	err = baseMultipartBuffer.BodyBufferWriter.Close()
	if err != nil {
		return err
	}

	// Prepare the request
	req, err := http.NewRequest("PUT", staticFilesEndpoint, baseMultipartBuffer.BodyBuffer)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", baseMultipartBuffer.BodyBufferWriter.FormDataContentType())

	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)

	// Forward error message if any
	microserviceError := sharedInfrastructure.ParseMicroserviceError(res, err)
	if microserviceError != nil {
		return microserviceError
	}

	return nil
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
	staticFilesEndpoint := fmt.Sprintf("%s/archives/delete", sharedInfrastructure.GetEnvironment().StaticFilesMicroserviceAddress)

	for _, archive := range archives {
		// Create the request body
		var body bytes.Buffer
		err := json.NewEncoder(&body).Encode(archive)
		if err != nil {
			errMessage := fmt.Sprintf("[ERR] - [BlocksPostgresRepository] - [deleteDependentArchives]: Unable to encode the request: %s", err.Error())
			log.Println(errMessage)
		}

		// Create the request
		req, err := http.NewRequest("POST", staticFilesEndpoint, &body)
		if err != nil {
			errMessage := fmt.Sprintf("[ERR] - [BlocksPostgresRepository] - [deleteDependentArchives]: Unable to create the request: %s", err.Error())
			log.Println(errMessage)
		}

		// Send the request
		client := &http.Client{}
		res, err := client.Do(req)

		// Forward error message if any
		microserviceError := sharedInfrastructure.ParseMicroserviceError(res, err)
		if microserviceError != nil {
			errMessage := fmt.Sprintf("[ERR] - [BlocksPostgresRepository] - [deleteDependentArchives]: Microservice returned an error: %s", microserviceError.Error())
			log.Println(errMessage)
		}
	}
}
