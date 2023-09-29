package application

import (
	"database/sql"

	"github.com/UPB-Code-Labs/main-api/src/courses/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/entities"
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/errors"
)

type CoursesUseCases struct {
	Repository              definitions.CoursesRepository
	InvitationCodeGenerator definitions.InvitationCodeGenerator
}

func (useCases *CoursesUseCases) GetRandomColor() (*entities.Color, error) {
	return useCases.Repository.GetRandomColor()
}

func (useCases *CoursesUseCases) GetInvitationCode(dto dtos.GetInvitationCodeDTO) (string, error) {
	// Check the teacher owns the course
	course, err := useCases.Repository.GetCourseByUUID(dto.CourseUUID)
	if err != nil {
		return "", err
	}

	teacherOwnsCourse := course.TeacherUUID == dto.TeacherUUID
	if !teacherOwnsCourse {
		return "", errors.TeacherDoesNotOwnsCourseError{}
	}

	// Return the code if it exists
	code, err := useCases.Repository.GetInvitationCode(dto.CourseUUID)
	unexpectedError := err != nil && err != sql.ErrNoRows
	if unexpectedError {
		return "", err
	}

	if code != "" {
		return code, nil
	}

	// Generate a new code if it does not exist
	code, err = useCases.generateInvitationCode(dto.CourseUUID)
	return code, err
}

func (useCases *CoursesUseCases) generateInvitationCode(courseUUID string) (string, error) {
	var generatedCode string

	// Generate a new code validating that it is unique
	generateNewCode := true
	for generateNewCode {
		code, err := useCases.InvitationCodeGenerator.Generate()
		if err != nil {
			return "", err
		}

		// Check if the code is unique
		course, err := useCases.Repository.GetCourseByInvitationCode(code)
		unexpectedError := err != nil && err != sql.ErrNoRows
		if unexpectedError {
			return "", err
		}

		isUnique := course == nil
		if isUnique {
			generateNewCode = false
			generatedCode = code
		}
	}

	// Save the code in the database
	err := useCases.Repository.SaveInvitationCode(courseUUID, generatedCode)
	if err != nil {
		return "", err
	}

	return generatedCode, nil
}

func (useCases *CoursesUseCases) SaveCourse(dto *dtos.CreateCourseDTO) (*entities.Course, error) {
	return useCases.Repository.SaveCourse(dto)
}

func (useCases *CoursesUseCases) JoinCourseUsingInvitationCode(dto *dtos.JoinCourseUsingInvitationCodeDTO) error {
	// Get the course by the invitation code
	course, err := useCases.Repository.GetCourseByInvitationCode(dto.InvitationCode)
	if err != nil {
		// Throw a domain error if no course with the given invitation code was found
		if err == sql.ErrNoRows {
			return errors.NoCourseWithInvitationCodeError{
				Code: dto.InvitationCode,
			}
		}

		return err
	}

	// Check if the student is already in the course
	isStudentInCourse, err := useCases.Repository.IsStudentInCourse(dto.StudentUUID, course.UUID)
	if err != nil {
		return err
	}
	if isStudentInCourse {
		return errors.StudentAlreadyInCourse{
			CourseName: course.Name,
		}
	}

	err = useCases.Repository.AddStudentToCourse(dto.StudentUUID, course.UUID)
	return err
}
