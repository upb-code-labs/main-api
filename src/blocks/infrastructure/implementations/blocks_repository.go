package implementations

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"

	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/errors"
	laboratoriesDomainErrors "github.com/UPB-Code-Labs/main-api/src/laboratories/domain/errors"
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
	baseMultipartBuffer, err := repository.getMultipartFormBuffer(
		staticFilesEndpoint,
		file,
	)
	if err != nil {
		return "", err
	}

	// Add the file type field to the request
	err = baseMultipartBuffer.bodyBufferWriter.WriteField("archive_type", "test")
	if err != nil {
		return "", err
	}

	// Close the writer
	err = baseMultipartBuffer.bodyBufferWriter.Close()
	if err != nil {
		return "", err
	}

	// Prepare the request
	req, err := http.NewRequest("POST", staticFilesEndpoint, baseMultipartBuffer.bodyBuffer)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", baseMultipartBuffer.bodyBufferWriter.FormDataContentType())

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
	baseMultipartBuffer, err := repository.getMultipartFormBuffer(
		staticFilesEndpoint,
		file,
	)
	if err != nil {
		return err
	}

	// Add the file type field to the request
	err = baseMultipartBuffer.bodyBufferWriter.WriteField("archive_type", "test")
	if err != nil {
		return err
	}

	// Add the archive uuid field to the request
	err = baseMultipartBuffer.bodyBufferWriter.WriteField("archive_uuid", uuid)
	if err != nil {
		return err
	}

	// Close the writer
	err = baseMultipartBuffer.bodyBufferWriter.Close()
	if err != nil {
		return err
	}

	// Prepare the request
	req, err := http.NewRequest("PUT", staticFilesEndpoint, baseMultipartBuffer.bodyBuffer)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", baseMultipartBuffer.bodyBufferWriter.FormDataContentType())

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

type baseMultipartFormBuffer struct {
	bodyBuffer       *bytes.Buffer
	bodyBufferWriter *multipart.Writer
}

func (repository *BlocksPostgresRepository) getMultipartFormBuffer(endpoint string, file *multipart.File) (br *baseMultipartFormBuffer, err error) {
	FILE_NAME := "archive.zip"
	FILE_CONTENT_TYPE := "application/zip"

	// Create multipart writer
	var bodyBuffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&bodyBuffer)

	// Add the file field to the request
	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition",
		fmt.Sprintf(
			`form-data; name="%s"; filename="%s"`,
			"file",
			FILE_NAME,
		),
	)
	header.Set("Content-Type", FILE_CONTENT_TYPE)

	// Reset the file pointer to the beginning
	_, err = (*file).Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Add the file to the request
	fileWriter, err := multipartWriter.CreatePart(header)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(fileWriter, *file); err != nil {
		return nil, err
	}

	return &baseMultipartFormBuffer{
		bodyBuffer:       &bodyBuffer,
		bodyBufferWriter: multipartWriter,
	}, nil
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
