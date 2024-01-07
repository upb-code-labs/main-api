package definitions

import (
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/entities"
)

type SubmissionsRepository interface {
	SaveSubmission(dto *dtos.CreateSubmissionDTO) (submissionUUID string, err error)
	ResetSubmissionStatus(submissionUUID string) (err error)

	GetStudentSubmission(studentUUID string, testBlockUUID string) (submission *entities.Submission, err error)
	GetSubmissionWorkMetadata(submissionUUID string) (submissionWorkMetadata *entities.SubmissionWork, err error)

	GetStudentSubmissionArchiveUUIDFromSubmissionUUID(submissionUUID string) (archiveUUID string, err error)
}
