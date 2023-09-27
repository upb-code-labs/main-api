package implementations

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"github.com/UPB-Code-Labs/main-api/src/courses/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/entities"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
)

type CoursesPostgresRepository struct {
	Connection *sql.DB
	Randomizer *rand.Rand
}

// Singleton
var coursesRepositoryInstance *CoursesPostgresRepository

func GetCoursesPgRepository() *CoursesPostgresRepository {
	if coursesRepositoryInstance == nil {
		randSource := rand.NewSource(time.Now().UnixNano())
		randInstance := rand.New(randSource)

		coursesRepositoryInstance = &CoursesPostgresRepository{
			Connection: infrastructure.GetPostgresConnection(),
			Randomizer: randInstance,
		}
	}

	return coursesRepositoryInstance
}

// Cache
var colorsLocaleStorage []entities.Color

// Methods
func (repository *CoursesPostgresRepository) GetRandomColor() (*entities.Color, error) {
	// Populate cache if empty
	if len(colorsLocaleStorage) == 0 {
		err := repository.populateColorsCache()
		if err != nil {
			return nil, err
		}
	}

	// Get random color
	if len(colorsLocaleStorage) > 0 {
		randomIndex := repository.Randomizer.Intn(len(colorsLocaleStorage))
		return &colorsLocaleStorage[randomIndex], nil
	} else {
		return nil, errors.New("no colors available")
	}
}

func (repository *CoursesPostgresRepository) populateColorsCache() error {
	repository.clearColorsStorage()

	// Query colors
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := "SELECT id, hexadecimal FROM colors"

	rows, err := repository.Connection.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Parse colors
	for rows.Next() {
		var color entities.Color

		err := rows.Scan(&color.UUID, &color.Hexadecimal)
		if err != nil {
			repository.clearColorsStorage()
			return err
		}

		colorsLocaleStorage = append(colorsLocaleStorage, color)
	}

	return nil
}

func (repository *CoursesPostgresRepository) clearColorsStorage() {
	colorsLocaleStorage = []entities.Color{}
}

func (repository *CoursesPostgresRepository) SaveCourse(dto *dtos.CreateCourseDTO) (*entities.Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Start transaction
	tx, err := repository.Connection.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create course
	createCourseQuery := `
		INSERT INTO courses (name, teacher_id, color_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	var courseId string
	row := tx.QueryRowContext(
		ctx,
		createCourseQuery,
		dto.Name,
		dto.TeacherUUID,
		dto.Color.UUID,
	)
	if row.Err() != nil {
		return nil, row.Err()
	}

	err = row.Scan(&courseId)
	if err != nil {
		return nil, err
	}

	// Add teacher to course
	addTeacherQuery := `
		INSERT INTO courses_has_users (class_id, user_id)
		VALUES ($1, $2)
	`

	_, err = tx.ExecContext(
		ctx,
		addTeacherQuery,
		courseId,
		dto.TeacherUUID,
	)
	if err != nil {
		return nil, err
	}

	// Commit changes
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return repository.GetCourseByUUID(courseId)
}

func (repository *CoursesPostgresRepository) GetCourseByUUID(uuid string) (*entities.Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := "SELECT id, teacher_id, name, color FROM courses_with_color WHERE id = $1"
	row := repository.Connection.QueryRowContext(ctx, query, uuid)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var course entities.Course
	err := row.Scan(&course.UUID, &course.TeacherUUID, &course.Name, &course.Color)
	if err != nil {
		return nil, err
	}

	return &course, nil
}

func (repository *CoursesPostgresRepository) GetCourseByInvitationCode(invitationCode string) (*entities.Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Get the course id
	getCourseIdQuery := "SELECT class_id FROM invitation_codes WHERE code = $1"
	row := repository.Connection.QueryRowContext(ctx, getCourseIdQuery, invitationCode)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var courseId string
	err := row.Scan(&courseId)
	if err != nil {
		return nil, err
	}

	// Get the course
	course, err := repository.GetCourseByUUID(courseId)
	if err != nil {
		return nil, err
	}

	return course, nil
}
