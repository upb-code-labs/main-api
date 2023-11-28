package dtos

type CreateLaboratoryDTO struct {
	TeacherUUID string
	CourseUUID  string
	Name        string
	OpeningDate string
	DueDate     string
}

type GetLaboratoryDTO struct {
	LaboratoryUUID string
	UserUUID       string
}

type UpdateLaboratoryDTO struct {
	LaboratoryUUID string
	TeacherUUID    string
	RubricUUID     *string
	Name           string
	OpeningDate    string
	DueDate        string
}
