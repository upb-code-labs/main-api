package definitions

import "github.com/UPB-Code-Labs/main-api/src/grades/domain/dtos"

// GradesRepository interface to be implemented by the repository
type GradesRepository interface {
	// GetStudentsGradesInLaboratory returns the grades of the students in a laboratory
	// that were graded using the given rubric by the teacher
	GetStudentsGradesInLaboratory(laboratoryUUID, rubricUUID string) ([]*dtos.SummarizedStudentGradeDTO, error)
}
