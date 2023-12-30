package dtos

import "github.com/UPB-Code-Labs/main-api/src/courses/domain/entities"

type AddStudentToCourseDTO struct {
	TeacherUUID string
	StudentUUID string
	CourseUUID  string
}

type CreateCourseDTO struct {
	Name        string
	TeacherUUID string
	Color       entities.Color
}

type EnrolledCoursesDto struct {
	Courses       []entities.Course
	HiddenCourses []entities.Course
}

type EnrolledStudentDTO struct {
	UUID            string
	FullName        string
	Email           string
	InstitutionalId string
	IsActive        bool
}

type SetUserStatusDTO struct {
	TeacherUUID string
	UserUUID    string
	CourseUUID  string
	ToActive    bool
}

type GetInvitationCodeDTO struct {
	CourseUUID  string
	TeacherUUID string
}

type JoinCourseUsingInvitationCodeDTO struct {
	StudentUUID    string
	InvitationCode string
}

type RenameCourseDTO struct {
	TeacherUUID string
	CourseUUID  string
	NewName     string
}
