package implementations

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"github.com/UPB-Code-Labs/main-api/src/courses/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/entities"
	courses_errors "github.com/UPB-Code-Labs/main-api/src/courses/domain/errors"
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
		INSERT INTO courses_has_users (course_id, user_id)
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

func (repository *CoursesPostgresRepository) SaveInvitationCode(courseUUID string, invitationCode string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		INSERT INTO invitation_codes (course_id, code)
		VALUES ($1, $2)
	`

	_, err := repository.Connection.ExecContext(
		ctx,
		query,
		courseUUID,
		invitationCode,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repository *CoursesPostgresRepository) GetInvitationCode(courseUUID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := "SELECT code FROM invitation_codes WHERE course_id = $1"
	row := repository.Connection.QueryRowContext(ctx, query, courseUUID)
	if row.Err() != nil {
		return "", row.Err()
	}

	var invitationCode string
	err := row.Scan(&invitationCode)
	if err != nil {
		return "", err
	}

	return invitationCode, nil
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
		// Throw a domain error if the course was not found
		if err == sql.ErrNoRows {
			return nil, courses_errors.NoCourseWithUUIDFound{
				UUID: uuid,
			}
		}

		return nil, err
	}

	return &course, nil
}

func (repository *CoursesPostgresRepository) GetCourseByInvitationCode(invitationCode string) (*entities.Course, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Get the course id
	getCourseIdQuery := "SELECT course_id FROM invitation_codes WHERE code = $1"
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

func (repository *CoursesPostgresRepository) AddStudentToCourse(studentUUID, courseUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		INSERT INTO courses_has_users (course_id, user_id)
		VALUES ($1, $2)
	`

	_, err := repository.Connection.ExecContext(
		ctx,
		query,
		courseUUID,
		studentUUID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repository *CoursesPostgresRepository) IsUserInCourse(userUUID, courseUUID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT COUNT(user_id) > 0
		FROM courses_has_users_view
		WHERE course_id = $1 AND 
		user_id = $2 AND 
		is_user_active = TRUE
	`

	row := repository.Connection.QueryRowContext(ctx, query, courseUUID, userUUID)
	if row.Err() != nil {
		return false, row.Err()
	}

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (repository *CoursesPostgresRepository) GetEnrolledCourses(studentUUID string) (*dtos.EnrolledCoursesDto, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT course_id, course_teacher_id, course_name, course_color, is_class_hidden
		FROM courses_has_users_view
		WHERE user_id = $1
		AND is_user_active = TRUE
	`

	rows, err := repository.Connection.QueryContext(ctx, query, studentUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	enrolledCourses := dtos.EnrolledCoursesDto{
		Courses:       []entities.Course{},
		HiddenCourses: []entities.Course{},
	}
	for rows.Next() {
		var course entities.Course
		var isClassHidden bool

		err := rows.Scan(
			&course.UUID,
			&course.TeacherUUID,
			&course.Name,
			&course.Color,
			&isClassHidden,
		)
		if err != nil {
			return nil, err
		}

		if isClassHidden {
			enrolledCourses.HiddenCourses = append(enrolledCourses.HiddenCourses, course)
		} else {
			enrolledCourses.Courses = append(enrolledCourses.Courses, course)
		}
	}

	return &enrolledCourses, nil
}

func (repository *CoursesPostgresRepository) ToggleCourseVisibility(courseUUID, studentUUID string) (isHiddenAfterUpdate bool, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		UPDATE courses_has_users
		SET is_class_hidden = NOT is_class_hidden
		WHERE course_id = $1 AND user_id = $2
		RETURNING is_class_hidden
	`

	row := repository.Connection.QueryRowContext(ctx, query, courseUUID, studentUUID)
	if row.Err() != nil {
		return false, row.Err()
	}

	err = row.Scan(&isHiddenAfterUpdate)
	return isHiddenAfterUpdate, err
}

func (repository *CoursesPostgresRepository) UpdateCourseName(dto dtos.RenameCourseDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := "UPDATE courses SET name = $1 WHERE id = $2"

	_, err := repository.Connection.ExecContext(
		ctx,
		query,
		dto.NewName,
		dto.CourseUUID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (repository *CoursesPostgresRepository) GetEnrolledStudents(courseUUID string) ([]*dtos.EnrolledStudentDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT user_id, user_full_name, user_email, user_institutional_id, is_user_active
		FROM courses_has_users_view
		WHERE course_id = $1 AND user_role = 'student'
	`

	rows, err := repository.Connection.QueryContext(ctx, query, courseUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	enrolledStudents := []*dtos.EnrolledStudentDTO{}
	for rows.Next() {
		var student dtos.EnrolledStudentDTO

		err := rows.Scan(
			&student.UUID,
			&student.FullName,
			&student.Email,
			&student.InstitutionalId,
			&student.IsActive,
		)
		if err != nil {
			return nil, err
		}

		enrolledStudents = append(enrolledStudents, &student)
	}

	return enrolledStudents, nil
}

func (repository *CoursesPostgresRepository) DoesTeacherOwnsCourse(teacherUUID, courseUUID string) (bool, error) {
	course, err := repository.GetCourseByUUID(courseUUID)
	if err != nil {
		return false, err
	}

	return course.TeacherUUID == teacherUUID, nil
}
