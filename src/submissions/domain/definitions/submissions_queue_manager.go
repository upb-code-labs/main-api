package definitions

import "github.com/UPB-Code-Labs/main-api/src/submissions/domain/entities"

type SubmissionsQueueManager interface {
	QueueWork(work *entities.SubmissionWork) (err error)
}
