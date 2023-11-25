package application

import (
	courses_definitions "github.com/UPB-Code-Labs/main-api/src/courses/domain/definitions"
	courses_errors "github.com/UPB-Code-Labs/main-api/src/courses/domain/errors"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/entities"
)

type LaboratoriesUseCases struct {
	LaboratoriesRepository definitions.LaboratoriesRepository
	CoursesRepository      courses_definitions.CoursesRepository
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

	return useCases.LaboratoriesRepository.SaveLaboratory(dto)
}
