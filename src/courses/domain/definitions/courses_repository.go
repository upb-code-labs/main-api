package definitions

import "github.com/UPB-Code-Labs/main-api/src/courses/domain/entities"

type CoursesRepository interface {
	SaveCourse(course *entities.Course) error

	GetCourseByUUID(uuid string) (*entities.Course, error)
	GetCourseByInvitationCode(invitationCode string) (*entities.Course, error)
}
