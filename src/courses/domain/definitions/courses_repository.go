package definitions

import (
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/entities"
)

type CoursesRepository interface {
	SaveCourse(dto *dtos.CreateCourseDTO) (*entities.Course, error)

	GetRandomColor() (*entities.Color, error)
	GetCourseByUUID(uuid string) (*entities.Course, error)
	GetCourseByInvitationCode(invitationCode string) (*entities.Course, error)
}
