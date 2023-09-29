package definitions

import (
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/entities"
)

type CoursesRepository interface {
	SaveCourse(dto *dtos.CreateCourseDTO) (*entities.Course, error)
	SaveInvitationCode(courseUUID string, invitationCode string) error
	GetInvitationCode(courseUUID string) (string, error)
	GetCourseByUUID(uuid string) (*entities.Course, error)
	GetCourseByInvitationCode(invitationCode string) (*entities.Course, error)
	GetRandomColor() (*entities.Color, error)
}
