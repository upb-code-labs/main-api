package definitions

import (
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/entities"
)

type SubmissionsRepository interface {
	SaveSubmission(dto *dtos.CreateSubmissionDTO) (submissionUUID string, err error)
	GetSubmission(dto *dtos.GetSubmissionDTO) (submissions *entities.Submission, err error)
}
