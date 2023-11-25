package requests

import "github.com/UPB-Code-Labs/main-api/src/laboratories/domain/dtos"

type CreateLaboratoryRequest struct {
	CourseUUID  string `json:"course_uuid" validate:"required,uuid4"`
	Name        string `json:"name" validate:"required,min=4,max=255"`
	OpeningDate string `json:"opening_date" validate:"required,ISO_date"`
	DueDate     string `json:"due_date" validate:"required,ISO_date"`
}

func (request *CreateLaboratoryRequest) ToDTO(teacherUUID string) *dtos.CreateLaboratoryDTO {
	return &dtos.CreateLaboratoryDTO{
		TeacherUUID: teacherUUID,
		CourseUUID:  request.CourseUUID,
		Name:        request.Name,
		OpeningDate: request.OpeningDate,
		DueDate:     request.DueDate,
	}
}
