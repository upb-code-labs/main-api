package application

import "github.com/UPB-Code-Labs/main-api/src/submissions/domain/definitions"

type SubmissionUseCases struct {
	SubmissionsRepository   definitions.SubmissionsRepository
	SubmissionsQueueManager definitions.SubmissionsQueueManager
}
