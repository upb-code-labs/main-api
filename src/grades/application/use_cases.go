package application

import (
	gradesDefinitions "github.com/UPB-Code-Labs/main-api/src/grades/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/grades/domain/dtos"
	gradesErrors "github.com/UPB-Code-Labs/main-api/src/grades/domain/errors"
	laboratoriesDefinitions "github.com/UPB-Code-Labs/main-api/src/laboratories/domain/definitions"
	laboratoriesErrors "github.com/UPB-Code-Labs/main-api/src/laboratories/domain/errors"
	rubricsDefinitions "github.com/UPB-Code-Labs/main-api/src/rubrics/domain/definitions"
	rubricsErrors "github.com/UPB-Code-Labs/main-api/src/rubrics/domain/errors"
)

type GradesUseCases struct {
	GradesRepository       gradesDefinitions.GradesRepository
	LaboratoriesRepository laboratoriesDefinitions.LaboratoriesRepository
	RubricsRepository      rubricsDefinitions.RubricsRepository
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

// SetCriteriaToGrade sets a criteria to a student's grade
func (useCases *GradesUseCases) SetCriteriaToGrade(dto *dtos.SetCriteriaToGradeDTO) error {
	// Validate the teacher owns the laboratory
	teacherOwnsLaboratory, err := useCases.LaboratoriesRepository.DoesTeacherOwnLaboratory(
		dto.TeacherUUID,
		dto.LaboratoryUUID,
	)
	if err != nil {
		return err
	}
	if !teacherOwnsLaboratory {
		return laboratoriesErrors.TeacherDoesNotOwnLaboratoryError{}
	}

	// Get the UUID of the current rubric of the laboratory
	laboratoryUUID := dto.LaboratoryUUID
	laboratoryInformation, err := useCases.LaboratoriesRepository.GetLaboratoryInformationByUUID(laboratoryUUID)
	if err != nil {
		return err
	}

	// Return an error if the laboratory does not have a rubric
	rubricUUID := laboratoryInformation.RubricUUID
	if rubricUUID == nil {
		return gradesErrors.LaboratoryDoesNotHaveRubricError{}
	}

	// Set the rubric UUID
	dto.RubricUUID = *rubricUUID

	// Validate the objective belongs to the rubric
	objectiveBelongsToRubric, err := useCases.RubricsRepository.DoesRubricHaveObjective(
		*rubricUUID,
		dto.ObjectiveUUID,
	)
	if err != nil {
		return err
	}
	if !objectiveBelongsToRubric {
		return &rubricsErrors.ObjectiveDoesNotBelongToRubricError{}
	}

	// Validate the criteria belongs to the objective
	if dto.CriteriaUUID != nil {
		criteriaBelongsToObjective, err := useCases.RubricsRepository.DoesObjectiveHaveCriteria(
			dto.ObjectiveUUID,
			*dto.CriteriaUUID,
		)
		if err != nil {
			return err
		}
		if !criteriaBelongsToObjective {
			return &rubricsErrors.CriteriaDoesNotBelongToObjectiveError{}
		}
	}

	// Set the criteria to the student's grade
	return useCases.GradesRepository.SetCriteriaToGrade(dto)
}

// GetStudentGradeInLaboratoryWithRubric returns the grade of an student in a laboratory
// that was graded with an specific rubric
func (useCases *GradesUseCases) GetStudentGradeInLaboratoryWithRubric(dto *dtos.GetStudentGradeInLaboratoryWithRubricDTO) (*dtos.StudentGradeInLaboratoryWithRubricDTO, error) {
	// Check if the users ows the laboratory
	userOwnsLaboratory, err := useCases.LaboratoriesRepository.DoesTeacherOwnLaboratory(
		dto.UserUUID,
		dto.LaboratoryUUID,
	)
	if err != nil {
		return nil, err
	}

	// Check if the user is trying to get their own grade
	isSameStudent := dto.StudentUUID == dto.UserUUID

	// Validate the user owns the laboratory or is the student
	userHasPermission := userOwnsLaboratory || isSameStudent
	if !userHasPermission {
		return nil, gradesErrors.UserCannotReadGradeError{}
	}

	// Get the grade
	grade, err := useCases.GradesRepository.GetStudentGradeInLaboratoryWithRubric(dto)
	return grade, err
}

// SetCommentToGrade sets a comment to an student's grade
func (useCases *GradesUseCases) SetCommentToGrade(dto *dtos.SetCommentToGradeDTO) error {
	// Validate the teacher owns the laboratory
	teacherOwnsLaboratory, err := useCases.LaboratoriesRepository.DoesTeacherOwnLaboratory(
		dto.TeacherUUID,
		dto.LaboratoryUUID,
	)
	if err != nil {
		return err
	}
	if !teacherOwnsLaboratory {
		return laboratoriesErrors.TeacherDoesNotOwnLaboratoryError{}
	}

	// Get the UUID of the current rubric of the laboratory
	laboratoryUUID := dto.LaboratoryUUID
	laboratoryInformation, err := useCases.LaboratoriesRepository.GetLaboratoryInformationByUUID(laboratoryUUID)
	if err != nil {
		return err
	}

	// Return an error if the laboratory does not have a rubric
	rubricUUID := laboratoryInformation.RubricUUID
	if rubricUUID == nil {
		return gradesErrors.LaboratoryDoesNotHaveRubricError{}
	}

	// Set the rubric UUID
	dto.RubricUUID = *rubricUUID

	// Set the comment to the student's grade
	return useCases.GradesRepository.SetCommentToGrade(dto)
}
