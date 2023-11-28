package dtos

type CreateLaboratoryDTO struct {
	TeacherUUID string
	CourseUUID  string
	Name        string
	OpeningDate string
	DueDate     string
}

type UpdateLaboratoryDTO struct {
	LaboratoryUUID string
	TeacherUUID    string
	RubricUUID     *string
	Name           string
	OpeningDate    string
	DueDate        string
}
