package definitions

import "github.com/UPB-Code-Labs/main-api/src/grades/domain/dtos"

// GradesRepository interface to be implemented by the repository
type GradesRepository interface {
	GetStudentsGradesInLaboratory(
		laboratoryUUID,
		rubricUUID string,
	) (
		[]*dtos.SummarizedStudentGradeDTO, error,
	)
	SetCriteriaToGrade(dto *dtos.SetCriteriaToGradeDTO) error
	SetCommentToGrade(dto *dtos.SetCommentToGradeDTO) error
	GetStudentGradeInLaboratoryWithRubric(
		dto *dtos.GetStudentGradeInLaboratoryWithRubricDTO,
	) (
		*dtos.StudentGradeInLaboratoryWithRubricDTO,
		error,
	)
}
