package implementations

import (
	"database/sql"

	"github.com/UPB-Code-Labs/main-api/src/grades/domain/dtos"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
)

// GradesPostgresRepository implementation of the GradesRepository interface
type GradesPostgresRepository struct {
	Connection *sql.DB
}

var gradesRepositoryInstance *GradesPostgresRepository

// GetGradesPostgresRepositoryInstance returns the singleton instance of the GradesPostgresRepository
func GetGradesPostgresRepositoryInstance() *GradesPostgresRepository {
	if gradesRepositoryInstance == nil {
		gradesRepositoryInstance = &GradesPostgresRepository{
			Connection: sharedInfrastructure.GetPostgresConnection(),
		}
	}

	return gradesRepositoryInstance
}

// GetStudentsGradesInLaboratory returns the grades of the students in a laboratory
// that were graded using the current rubric of the laboratory by the teacher
func (repository *GradesPostgresRepository) GetStudentsGradesInLaboratory(laboratoryUUID string) ([]*dtos.SummarizedStudentGradeDTO, error) {
	return nil, nil
}
