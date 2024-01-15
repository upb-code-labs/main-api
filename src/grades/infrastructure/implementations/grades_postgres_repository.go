package implementations

import (
	"context"
	"database/sql"
	"time"

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
func (repository *GradesPostgresRepository) GetStudentsGradesInLaboratory(laboratoryUUID, rubricUUID string) ([]*dtos.SummarizedStudentGradeDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	query := `
		SELECT g.student_id, u.full_name, g.grade
		FROM grades g
		INNER JOIN users u ON u.id = g.student_id
		WHERE g.laboratory_id = $1 AND g.rubric_id = $2
	`

	// Run the query
	rows, err := repository.Connection.QueryContext(ctx, query, laboratoryUUID, rubricUUID)
	if err != nil {
		return nil, err
	}

	// Parse the results
	var summarizedGrades []*dtos.SummarizedStudentGradeDTO
	for rows.Next() {
		var studentGrade dtos.SummarizedStudentGradeDTO

		if err := rows.Scan(
			&studentGrade.StudentUUID,
			&studentGrade.StudentFullName,
			&studentGrade.Grade); err != nil {
			return nil, err
		}

		summarizedGrades = append(summarizedGrades, &studentGrade)
	}

	return summarizedGrades, nil
}
