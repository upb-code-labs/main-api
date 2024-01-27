package requests

import (
	"time"

	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/dtos"
)

type CreateLaboratoryRequest struct {
	CourseUUID  string `json:"course_uuid" validate:"required,uuid4"`
	Name        string `json:"name" validate:"required,min=4,max=255"`
	OpeningDate string `json:"opening_date" validate:"required,RFC3339_date"`
	DueDate     string `json:"due_date" validate:"required,RFC3339_date"`
}

func (request *CreateLaboratoryRequest) ToDTO(teacherUUID string) *dtos.CreateLaboratoryDTO {
	parsedOpeningDate, _ := time.Parse(time.RFC3339, request.OpeningDate)
	parsedDueDate, _ := time.Parse(time.RFC3339, request.DueDate)

	return &dtos.CreateLaboratoryDTO{
		TeacherUUID: teacherUUID,
		CourseUUID:  request.CourseUUID,
		Name:        request.Name,
		OpeningDate: parsedOpeningDate,
		DueDate:     parsedDueDate,
	}
}

type UpdateLaboratoryRequest struct {
	Name        string  `json:"name" validate:"required,min=4,max=255"`
	OpeningDate string  `json:"opening_date" validate:"required,RFC3339_date"`
	DueDate     string  `json:"due_date" validate:"required,RFC3339_date"`
	RubricUUID  *string `json:"rubric_uuid,omitempty" validate:"omitempty,uuid4"`
}

func (request *UpdateLaboratoryRequest) ToDTO(laboratoryUUID string, teacherUUID string) *dtos.UpdateLaboratoryDTO {
	parsedOpeningDate, _ := time.Parse(time.RFC3339, request.OpeningDate)
	parsedDueDate, _ := time.Parse(time.RFC3339, request.DueDate)

	return &dtos.UpdateLaboratoryDTO{
		TeacherUUID:    teacherUUID,
		LaboratoryUUID: laboratoryUUID,
		RubricUUID:     request.RubricUUID,
		Name:           request.Name,
		OpeningDate:    parsedOpeningDate,
		DueDate:        parsedDueDate,
	}
}

type CreateTestBlockRequest struct {
	LaboratoryUUID string `validate:"required,uuid4"`
	LanguageUUID   string `validate:"required,uuid4"`
	Name           string `validate:"required,min=4,max=255"`
}
