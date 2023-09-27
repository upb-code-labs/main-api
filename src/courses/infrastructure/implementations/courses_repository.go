package implementations

import (
	"database/sql"

	"github.com/UPB-Code-Labs/main-api/src/courses/domain/entities"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
)

type CoursesPostgresRepository struct {
	Connection *sql.DB
}

var coursesRepositoryInstance *CoursesPostgresRepository

func GetCoursesPgRepository() *CoursesPostgresRepository {
	if coursesRepositoryInstance == nil {
		coursesRepositoryInstance = &CoursesPostgresRepository{
			Connection: infrastructure.GetPostgresConnection(),
		}
	}

	return coursesRepositoryInstance
}

func (repository *CoursesPostgresRepository) SaveCourse(course *entities.Course) error {
	return nil
}

func (repository *CoursesPostgresRepository) GetCourseByUUID(uuid string) (*entities.Course, error) {
	return nil, nil
}

func (repository *CoursesPostgresRepository) GetCourseByInvitationCode(invitationCode string) (*entities.Course, error) {
	return nil, nil
}
