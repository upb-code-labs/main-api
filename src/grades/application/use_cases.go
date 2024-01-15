package application

import (
	gradesDefinitions "github.com/UPB-Code-Labs/main-api/src/grades/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/grades/domain/dtos"
	laboratoriesDefinitions "github.com/UPB-Code-Labs/main-api/src/laboratories/domain/definitions"
	laboratoriesErrors "github.com/UPB-Code-Labs/main-api/src/laboratories/domain/errors"
)

type GradesUseCases struct {
	GradesRepository       gradesDefinitions.GradesRepository
	LaboratoriesRepository laboratoriesDefinitions.LaboratoriesRepository
}

// GetSummarizedGradesInLaboratory returns the summarized version (Just student's UUID, full name and grade) of the grades
func (useCases *GradesUseCases) GetSummarizedGradesInLaboratory(dto *dtos.GetSummarizedGradesInLaboratoryDTO) ([]*dtos.SummarizedStudentGradeDTO, error) {
	// Validate the teacher owns the laboratory
	teacherOwnsLaboratory, err := useCases.LaboratoriesRepository.DoesTeacherOwnLaboratory(
		dto.TeacherUUID,
		dto.LaboratoryUUID,
	)
	if err != nil {
		return nil, err
	}
	if !teacherOwnsLaboratory {
		return nil, laboratoriesErrors.TeacherDoesNotOwnLaboratoryError{}
	}

	// Get the UUID of the current rubric of the laboratory
	laboratoryUUID := dto.LaboratoryUUID
	laboratoryInformation, err := useCases.LaboratoriesRepository.GetLaboratoryInformationByUUID(laboratoryUUID)
	if err != nil {
		return nil, err
	}

	// Return an empty array if the laboratory does not have a rubric
	rubricUUID := laboratoryInformation.RubricUUID
	if rubricUUID == nil {
		return []*dtos.SummarizedStudentGradeDTO{}, nil
	}

	// Get the grades
	return useCases.GradesRepository.GetStudentsGradesInLaboratory(laboratoryUUID, *rubricUUID)
}
