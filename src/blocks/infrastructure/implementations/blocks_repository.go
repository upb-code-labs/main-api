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

func (repository *BlocksPostgresRepository) SaveTestsArchive(file *multipart.File) (uuid string, err error) {
	// Arbitrary constants since they are not used anywhere else and
	// we only support zip files for now.
	FILE_NAME := "archive.zip"
	FILE_CONTENT_TYPE := "application/zip"

	// Create multipart writer
	staticFilesEndpoint := fmt.Sprintf("%s/archives/save", sharedInfrastructure.GetEnvironment().StaticFilesMicroserviceAddress)
	var requestBuffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBuffer)

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
		return "", err
	}

	fileWriter, err := multipartWriter.CreatePart(header)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(fileWriter, *file); err != nil {
		return "", err
	}

	// Add the file type field to the request
	if err := multipartWriter.WriteField("archive_type", "test"); err != nil {
		return "", err
	}

	// Close the writer
	err = multipartWriter.Close()
	if err != nil {
		return "", err
	}

	// Prepare the request
	req, err := http.NewRequest("POST", staticFilesEndpoint, &requestBuffer)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)

	// Forward error message if any
	microserviceError := sharedInfrastructure.ParseMicroserviceError(res, err)
	if microserviceError != nil {
		return "", microserviceError
	}

	// Parse the response
	var response map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	if response["uuid"] == nil {
		return "", &sharedDomainErrors.StaticFilesMicroserviceError{
			Code:    http.StatusInternalServerError,
			Message: "The static files microservice did not return the UUID of the saved file",
		}
	}

	return response["uuid"].(string), nil
}
