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
	teacherOwnsCourse, err := useCases.Repository.DoesTeacherOwnsCourse(dto.TeacherUUID, dto.CourseUUID)
	if err != nil {
		return "", err
	}
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

func (useCases *CoursesUseCases) GetCourse(userUUID, courseUUID string) (*entities.Course, error) {
	// Get the course
	course, err := useCases.Repository.GetCourseByUUID(courseUUID)
	if err != nil {
		return nil, err
	}

	// Check the user is enrolled in the course
	isStudentInCourse, err := useCases.Repository.IsUserInCourse(userUUID, courseUUID)
	if err != nil {
		return nil, err
	}
	if !isStudentInCourse {
		return nil, errors.UserNotInCourseError{}
	}

	// Return the data
	return course, nil
}

func (useCases *CoursesUseCases) JoinCourseUsingInvitationCode(dto *dtos.JoinCourseUsingInvitationCodeDTO) (*entities.Course, error) {
	// Get the course by the invitation code
	course, err := useCases.Repository.GetCourseByInvitationCode(dto.InvitationCode)
	if err != nil {
		// Throw a domain error if no course with the given invitation code was found
		if err == sql.ErrNoRows {
			return nil, errors.NoCourseWithInvitationCodeError{
				Code: dto.InvitationCode,
			}
		}

		return nil, err
	}

	// Check if the student is already in the course
	isStudentInCourse, err := useCases.Repository.IsUserInCourse(dto.StudentUUID, course.UUID)
	if err != nil {
		return nil, err
	}
	if isStudentInCourse {
		return nil, errors.StudentAlreadyInCourse{
			CourseName: course.Name,
		}
	}

	err = useCases.Repository.AddStudentToCourse(dto.StudentUUID, course.UUID)
	return course, err
}

func (useCases *CoursesUseCases) GetEnrolledCourses(userUUID string) (*dtos.EnrolledCoursesDto, error) {
	return useCases.Repository.GetEnrolledCourses(userUUID)
}

func (useCases *CoursesUseCases) ToggleCourseVisibility(courseUUID, userUUID string) (bool, error) {
	// Check the user is enrolled in the course
	isStudentInCourse, err := useCases.Repository.IsUserInCourse(userUUID, courseUUID)
	if err != nil {
		return false, err
	}
	if !isStudentInCourse {
		return false, errors.UserNotInCourseError{}
	}

	// Toggle the course visibility for the user
	return useCases.Repository.ToggleCourseVisibility(courseUUID, userUUID)
}

func (useCases *CoursesUseCases) UpdateCourseName(dto dtos.RenameCourseDTO) error {
	// Check the user is the teacher of the course
	course, err := useCases.Repository.GetCourseByUUID(dto.CourseUUID)
	if err != nil {
		return err
	}

	teacherOwnsCourse := course.TeacherUUID == dto.TeacherUUID
	if !teacherOwnsCourse {
		return errors.TeacherDoesNotOwnsCourseError{}
	}

	// Check the new name is different from the current one
	if course.Name == dto.NewName {
		return errors.UnchangedCourseNameError{}
	}

	// Update the course name
	return useCases.Repository.UpdateCourseName(dto)
}

func (useCases *CoursesUseCases) AddStudentToCourse(dto *dtos.AddStudentToCourseDTO) error {
	// Check the user is the teacher of the course
	course, err := useCases.Repository.GetCourseByUUID(dto.CourseUUID)
	if err != nil {
		return err
	}

	teacherOwnsCourse := course.TeacherUUID == dto.TeacherUUID
	if !teacherOwnsCourse {
		return errors.TeacherDoesNotOwnsCourseError{}
	}

	// Check the student is not already in the course
	isStudentInCourse, err := useCases.Repository.IsUserInCourse(dto.StudentUUID, dto.CourseUUID)
	if err != nil {
		return err
	}
	if isStudentInCourse {
		return errors.StudentAlreadyInCourse{
			CourseName: course.Name,
		}
	}

	// Add the student to the course
	return useCases.Repository.AddStudentToCourse(dto.StudentUUID, dto.CourseUUID)
}

func (useCases *CoursesUseCases) GetEnrolledStudents(teacherUUID, courseUUID string) ([]*dtos.EnrolledStudentDTO, error) {
	// Check the teacher owns the course
	teacherOwnsCourse, err := useCases.Repository.DoesTeacherOwnsCourse(teacherUUID, courseUUID)
	if err != nil {
		return nil, err
	}
	if !teacherOwnsCourse {
		return nil, errors.TeacherDoesNotOwnsCourseError{}
	}

	// Get the enrolled students
	return useCases.Repository.GetEnrolledStudents(courseUUID)
}

func (useCases *CoursesUseCases) SetStudentStatus(dto *dtos.SetUserStatusDTO) error {
	// Get the course
	course, err := useCases.Repository.GetCourseByUUID(dto.CourseUUID)
	if err != nil {
		return err
	}

	// Check the user is the teacher of the course
	wantsToUpdateCourseTeacher := dto.UserUUID == course.TeacherUUID
	if wantsToUpdateCourseTeacher {
		return errors.CannotUpdateCourseTeacherStatus{}
	}

	teacherOwnsCourse := course.TeacherUUID == dto.TeacherUUID
	if !teacherOwnsCourse {
		return errors.TeacherDoesNotOwnsCourseError{}
	}

	// Update the student status
	return useCases.Repository.SetStudentStatus(dto)
}

func (useCases *CoursesUseCases) GetCourseLaboratories(dto dtos.GetCourseLaboratoriesDTO) ([]*dtos.BaseLaboratoryDTO, error) {
	// Check the user is enrolled in the course
	isUserInCourse, err := useCases.Repository.IsUserInCourse(dto.UserUUID, dto.CourseUUID)
	if err != nil {
		return nil, err
	}
	if !isUserInCourse {
		return nil, errors.UserNotInCourseError{}
	}

	// Get the course laboratories according to the user role
	if dto.UserRole == "teacher" {
		return useCases.Repository.GetCourseLaboratories(dto.CourseUUID)
	} else {
		return useCases.Repository.GetCourseActiveLaboratories(dto.CourseUUID)
	}
}
