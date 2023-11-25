package dtos

type CreateLaboratoryDTO struct {
	TeacherUUID string
	CourseUUID  string
	Name        string
	OpeningDate string
	DueDate     string
}
