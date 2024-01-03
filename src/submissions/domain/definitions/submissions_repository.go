package definitions

import (
	"mime/multipart"

	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/entities"
)

type SubmissionsRepository interface {
	// Methods to interact with the static files microservice
	SaveSubmissionArchive(file *multipart.File) (archiveUUID string, err error)
	OverwriteSubmissionArchive(file *multipart.File, archiveUUID string) (err error)

	SaveSubmission(dto *dtos.CreateSubmissionDTO) (submissionUUID string, err error)
	ResetSubmissionStatus(submissionUUID string) (err error)

	GetSubmission(dto *dtos.GetSubmissionDTO) (submissions *entities.Submission, err error)
	GetStudentSubmission(studentUUID string, testBlockUUID string) (submission *entities.Submission, err error)
	GetSubmissionWorkMetadata(submissionUUID string) (submissionWorkMetadata *entities.SubmissionWork, err error)

	// Get the UUID of the .zip archive saved in the static files microservice for a given submission
	GetStudentSubmissionArchiveUUIDFromSubmissionUUID(submissionUUID string) (archiveUUID string, err error)
}
