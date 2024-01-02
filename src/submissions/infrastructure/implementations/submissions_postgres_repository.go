package implementations

import (
	"database/sql"

	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/submissions/domain/entities"
)

type SubmissionsPostgresRepository struct {
	Connection *sql.DB
}

// Singleton
var submissionsRepositoryInstance *SubmissionsPostgresRepository

func GetSubmissionsRepositoryInstance() *SubmissionsPostgresRepository {
	if submissionsRepositoryInstance == nil {
		submissionsRepositoryInstance = &SubmissionsPostgresRepository{
			Connection: sharedInfrastructure.GetPostgresConnection(),
		}
	}

	return submissionsRepositoryInstance
}

// Methods implementation
func (repository *SubmissionsPostgresRepository) SaveSubmission(dto *dtos.CreateSubmissionDTO) (submissionUUID string, err error) {
	return "", nil
}

func (repository *SubmissionsPostgresRepository) GetSubmission(dto *dtos.GetSubmissionDTO) (submission *entities.Submission, err error) {
	return nil, nil
}
