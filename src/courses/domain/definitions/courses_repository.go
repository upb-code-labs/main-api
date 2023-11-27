package definitions

import (
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/entities"
)

type CoursesRepository interface {
	SaveCourse(dto *dtos.CreateCourseDTO) (*entities.Course, error)
	GetCourseByUUID(uuid string) (*entities.Course, error)

	SaveInvitationCode(courseUUID, invitationCode string) error
	GetInvitationCode(courseUUID string) (string, error)
	GetCourseByInvitationCode(invitationCode string) (*entities.Course, error)

	AddStudentToCourse(studentUUID, courseUUID string) error
	IsUserInCourse(userUUID, courseUUID string) (bool, error)
	GetEnrolledCourses(studentUUID string) (*dtos.EnrolledCoursesDto, error)
	GetEnrolledStudents(courseUUID string) ([]*dtos.EnrolledStudentDTO, error)

	GetRandomColor() (*entities.Color, error)
	ToggleCourseVisibility(courseUUID, studentUUID string) (isHiddenAfterUpdate bool, err error)
	UpdateCourseName(dtos.RenameCourseDTO) error

	DoesTeacherOwnsCourse(teacherUUID, courseUUID string) (bool, error)
}
