package errors

import (
	"fmt"
	"net/http"
)

type NoCourseWithInvitationCodeError struct {
	Code string
}

func (err NoCourseWithInvitationCodeError) Error() string {
	return fmt.Sprintf("Course with invitation code %s not found", err.Code)
}

func (err NoCourseWithInvitationCodeError) StatusCode() int {
	return http.StatusNotFound
}

type NoCourseWithUUIDFound struct {
	UUID string
}

func (err NoCourseWithUUIDFound) Error() string {
	return fmt.Sprintf("Course with UUID %s not found", err.UUID)
}

func (err NoCourseWithUUIDFound) StatusCode() int {
	return http.StatusNotFound
}

type StudentAlreadyInCourse struct {
	CourseName string
}

func (err StudentAlreadyInCourse) Error() string {
	return fmt.Sprintf("Student is already in the course %s", err.CourseName)
}

func (err StudentAlreadyInCourse) StatusCode() int {
	return http.StatusConflict
}

type TeacherDoesNotOwnsCourseError struct {
}

func (err TeacherDoesNotOwnsCourseError) Error() string {
	return "You do not own the course"
}

func (err TeacherDoesNotOwnsCourseError) StatusCode() int {
	return http.StatusForbidden
}

type UnchangedCourseNameError struct {
}

func (err UnchangedCourseNameError) Error() string {
	return "The course has the same name"
}

func (err UnchangedCourseNameError) StatusCode() int {
	return http.StatusBadRequest
}

type UserNotInCourseError struct{}

func (err UserNotInCourseError) Error() string {
	return "You are not enrolled in the course"
}

func (err UserNotInCourseError) StatusCode() int {
	return http.StatusForbidden
}

type CannotUpdateCourseTeacherStatus struct{}

func (err CannotUpdateCourseTeacherStatus) Error() string {
	return "You cannot update the teacher status"
}

func (err CannotUpdateCourseTeacherStatus) StatusCode() int {
	return http.StatusConflict
}
