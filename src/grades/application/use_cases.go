package application

import (
	"github.com/UPB-Code-Labs/main-api/src/grades/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/grades/domain/dtos"
)

type GradesUseCases struct {
	GradesRepository definitions.GradesRepository
}

// GetSummarizedGradesInLaboratory returns the summarized version (Just student's UUID, full name and grade) of the grades
func (useCases *GradesUseCases) GetSummarizedGradesInLaboratory(dto *dtos.GetSummarizedGradesInLaboratoryDTO) ([]*dtos.SummarizedStudentGradeDTO, error) {
	return useCases.GradesRepository.GetStudentsGradesInLaboratory(dto.LaboratoryUUID)
}
