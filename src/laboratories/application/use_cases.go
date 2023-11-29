package application

import (
	courses_definitions "github.com/UPB-Code-Labs/main-api/src/courses/domain/definitions"
	courses_errors "github.com/UPB-Code-Labs/main-api/src/courses/domain/errors"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/entities"
	rubrics_definitions "github.com/UPB-Code-Labs/main-api/src/rubrics/domain/definitions"
	rubrics_errors "github.com/UPB-Code-Labs/main-api/src/rubrics/domain/errors"
)

type LaboratoriesUseCases struct {
	LaboratoriesRepository definitions.LaboratoriesRepository
	CoursesRepository      courses_definitions.CoursesRepository
	RubricsRepository      rubrics_definitions.RubricsRepository
}

func (useCases *LaboratoriesUseCases) CreateLaboratory(dto *dtos.CreateLaboratoryDTO) (laboratory *entities.Laboratory, err error) {
	// Check that the teacher owns the course
	ownsCourse, err := useCases.CoursesRepository.DoesTeacherOwnsCourse(dto.TeacherUUID, dto.CourseUUID)
	if err != nil {
		return nil, err
	}

	if !ownsCourse {
		return nil, courses_errors.TeacherDoesNotOwnsCourseError{}
	}

	// Create the laboratory
	return useCases.LaboratoriesRepository.SaveLaboratory(dto)
}

func (useCases *LaboratoriesUseCases) GetLaboratory(dto *dtos.GetLaboratoryDTO) (laboratory *entities.Laboratory, err error) {
	// Get the laboratory
	laboratory, err = useCases.LaboratoriesRepository.GetLaboratoryByUUID(dto.LaboratoryUUID)
	if err != nil {
		return nil, err
	}

	// Check that the user is enrolled in the course
	isEnrolled, err := useCases.CoursesRepository.IsUserInCourse(dto.UserUUID, laboratory.CourseUUID)
	if err != nil {
		return nil, err
	}

	if !isEnrolled {
		return nil, courses_errors.UserNotInCourseError{}
	}

	return laboratory, nil
}

func (useCases *LaboratoriesUseCases) UpdateLaboratory(dto *dtos.UpdateLaboratoryDTO) error {
	// Check that the teacher owns the laboratory / course
	teacherOwnsLaboratory, err := useCases.doesTeacherOwnsLaboratory(dto.TeacherUUID, dto.LaboratoryUUID)
	if err != nil {
		return err
	}

	if !teacherOwnsLaboratory {
		return &courses_errors.TeacherDoesNotOwnsCourseError{}
	}

	// Check that the teacher owns the rubric
	if dto.RubricUUID != nil {
		teacherOwnsRubric, err := useCases.RubricsRepository.DoesTeacherOwnRubric(dto.TeacherUUID, *dto.RubricUUID)
		if err != nil {
			return err
		}
		if !teacherOwnsRubric {
			return &rubrics_errors.TeacherDoesNotOwnsRubric{}
		}
	}

	// Update the laboratory
	return useCases.LaboratoriesRepository.UpdateLaboratory(dto)
}

func (useCases *LaboratoriesUseCases) CreateMarkdownBlock(dto *dtos.CreateMarkdownBlockDTO) (blockUUID string, err error) {
	// Check that the teacher owns the laboratory
	teacherOwnsLaboratory, err := useCases.doesTeacherOwnsLaboratory(dto.TeacherUUID, dto.LaboratoryUUID)
	if err != nil {
		return "", err
	}

	if !teacherOwnsLaboratory {
		return "", &courses_errors.TeacherDoesNotOwnsCourseError{}
	}

	// Create the block
	return useCases.LaboratoriesRepository.CreateMarkdownBlock(dto.LaboratoryUUID)
}

func (useCases *LaboratoriesUseCases) doesTeacherOwnsLaboratory(teacherUUID, laboratoryUUID string) (bool, error) {
	laboratory, err := useCases.LaboratoriesRepository.GetLaboratoryByUUID(laboratoryUUID)
	if err != nil {
		return false, err
	}

	return useCases.CoursesRepository.DoesTeacherOwnsCourse(teacherUUID, laboratory.CourseUUID)
}
