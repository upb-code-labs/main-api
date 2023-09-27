package application

import (
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/entities"
)

type CoursesUseCases struct {
	Repository              definitions.CoursesRepository
	InvitationCodeGenerator definitions.InvitationCodeGenerator
}

func (useCases *CoursesUseCases) GetRandomColor() (*entities.Color, error) {
	return useCases.Repository.GetRandomColor()
}

func (useCases *CoursesUseCases) SaveCourse(dto *dtos.CreateCourseDTO) (*entities.Course, error) {
	return useCases.Repository.SaveCourse(dto)
}
