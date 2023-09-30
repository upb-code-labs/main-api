package dtos

import "github.com/UPB-Code-Labs/main-api/src/courses/domain/entities"

type EnrolledCoursesDto struct {
	Courses       []entities.Course
	HiddenCourses []entities.Course
}
