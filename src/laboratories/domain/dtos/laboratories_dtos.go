package dtos

import "mime/multipart"

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

type CreateMarkdownBlockDTO struct {
	TeacherUUID    string
	LaboratoryUUID string
}

type CreateTestBlockDTO struct {
	TeacherUUID     string
	LaboratoryUUID  string
	LanguageUUID    string
	TestArchiveUUID string
	Name            string
	MultipartFile   *multipart.File
}
