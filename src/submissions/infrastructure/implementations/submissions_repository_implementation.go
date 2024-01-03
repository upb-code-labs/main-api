package implementations

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	sharedDomainErrors "github.com/UPB-Code-Labs/main-api/src/shared/domain/errors"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/entities"
)

type SubmissionsRepositoryImpl struct {
	Connection *sql.DB
}

// Singleton
var submissionsRepositoryInstance *SubmissionsRepositoryImpl

func GetSubmissionsRepositoryInstance() *SubmissionsRepositoryImpl {
	if submissionsRepositoryInstance == nil {
		submissionsRepositoryInstance = &SubmissionsRepositoryImpl{
			Connection: sharedInfrastructure.GetPostgresConnection(),
		}
	}

	return submissionsRepositoryInstance
}

// Methods implementation
func (repository *SubmissionsRepositoryImpl) SaveSubmissionArchive(file *multipart.File) (archiveUUID string, err error) {
	staticFilesEndpoint := fmt.Sprintf("%s/archives/save", sharedInfrastructure.GetEnvironment().StaticFilesMicroserviceAddress)

	// Create multipart writer
	baseMultipartBuffer, err := sharedInfrastructure.GetMultipartFormBufferFromFile(file)
	if err != nil {
		return "", err
	}

	// Add the file type field to the request
	err = baseMultipartBuffer.BodyBufferWriter.WriteField("archive_type", "submission")
	if err != nil {
		return "", err
	}

	// Close the writer
	err = baseMultipartBuffer.BodyBufferWriter.Close()
	if err != nil {
		return "", err
	}

	// Create the request
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

func (repository *SubmissionsRepositoryImpl) OverwriteSubmissionArchive(file *multipart.File, archiveUUID string) (err error) {
	staticFilesEndpoint := fmt.Sprintf("%s/archives/overwrite", sharedInfrastructure.GetEnvironment().StaticFilesMicroserviceAddress)

	// Create multipart writer
	baseMultipartBuffer, err := sharedInfrastructure.GetMultipartFormBufferFromFile(file)
	if err != nil {
		return err
	}

	// Add the file type field to the request
	err = baseMultipartBuffer.BodyBufferWriter.WriteField("archive_type", "submission")
	if err != nil {
		return err
	}

	// Add the archive UUID field to the request
	err = baseMultipartBuffer.BodyBufferWriter.WriteField("archive_uuid", archiveUUID)
	if err != nil {
		return err
	}

	// Close the writer
	err = baseMultipartBuffer.BodyBufferWriter.Close()
	if err != nil {
		return err
	}

	// Create the request
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

func (repository *SubmissionsRepositoryImpl) SaveSubmission(dto *dtos.CreateSubmissionDTO) (submissionUUID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Start the transaction
	tx, err := repository.Connection.BeginTx(ctx, nil)
	if err != nil {
		return "", err
	}

	// Create an entry in the archives table
	var dbArchiveUUID string
	query := `
		INSERT INTO archives (file_id)
		VALUES ($1)
		RETURNING id
	`

	err = tx.QueryRowContext(ctx, query, dto.SavedArchiveUUID).Scan(&dbArchiveUUID)
	if err != nil {
		return "", err
	}

	// Create an entry in the submissions table
	var dbSubmissionUUID string
	query = `
		INSERT INTO submissions (student_id, test_block_id, archive_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err = tx.QueryRowContext(
		ctx, query, dto.StudentUUID, dto.TestBlockUUID, dbArchiveUUID,
	).Scan(&dbSubmissionUUID)
	if err != nil {
		return "", err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return "", err
	}

	return dbSubmissionUUID, nil
}

func (repository *SubmissionsRepositoryImpl) ResetSubmissionStatus(submissionUUID string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	DEFAULT_PASSING_VALUE := false
	DEFAULT_STATUS_VALUE := "pending"
	DEFAULT_STDOUT_VALUE := ""

	query := `
		UPDATE submissions
		SET
			passing = $1,
			status = $2,
			stdout = $3,
			submitted_at = CURRENT_TIMESTAMP
	`

	_, err = repository.Connection.ExecContext(
		ctx, query, DEFAULT_PASSING_VALUE, DEFAULT_STATUS_VALUE, DEFAULT_STDOUT_VALUE,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repository *SubmissionsRepositoryImpl) GetStudentSubmission(studentUUID string, testBlockUUID string) (submission *entities.Submission, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT id, archive_id, passing, status, stdout
		FROM submissions
		WHERE student_id = $1 AND test_block_id = $2
	`

	submission = &entities.Submission{}

	err = repository.Connection.QueryRowContext(
		ctx, query, studentUUID, testBlockUUID,
	).Scan(
		&submission.UUID, &submission.ArchiveUUID, &submission.Passing, &submission.Status, &submission.Stdout,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return submission, nil
}

func (repository *SubmissionsRepositoryImpl) GetSubmission(dto *dtos.GetSubmissionDTO) (submission *entities.Submission, err error) {
	return nil, nil
}

func (repository *SubmissionsRepositoryImpl) GetSubmissionWorkMetadata(submissionUUID string) (submissionWorkMetadata *entities.SubmissionWork, err error) {
	return nil, nil
}

func (repository *SubmissionsRepositoryImpl) GetStudentSubmissionArchiveUUIDFromSubmissionUUID(submissionUUID string) (archiveUUID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT file_id
		FROM archives
		WHERE id = (
			SELECT archive_id
			FROM submissions
			WHERE id = $1
		)
	`

	err = repository.Connection.QueryRowContext(
		ctx, query, submissionUUID,
	).Scan(&archiveUUID)

	if err != nil {
		return "", err
	}

	return archiveUUID, nil
}
