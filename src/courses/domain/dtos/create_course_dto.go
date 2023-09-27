package dtos

import "github.com/UPB-Code-Labs/main-api/src/courses/domain/entities"

type CreateCourseDTO struct {
	Name        string
	TeacherUUID string
	Color       entities.Color
}
