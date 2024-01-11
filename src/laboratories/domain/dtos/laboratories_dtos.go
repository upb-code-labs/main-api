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

type GetLaboratoryProgressDTO struct {
	LaboratoryUUID string
	TeacherUUID    string
}

type LaboratoryProgressDTO struct {
	TotalTestBlocks  int                             `json:"total_test_blocks"`
	StudentsProgress []*LaboratoryStudentProgressDTO `json:"students_progress"`
}

type LaboratoryStudentProgressDTO struct {
	StudentUUID        string `json:"student_uuid"`
	StudentFullName    string `json:"student_full_name"`
	PendingSubmissions int    `json:"pending_submissions"`
	RunningSubmissions int    `json:"running_submissions"`
	FailingSubmissions int    `json:"failing_submissions"`
	SuccessSubmissions int    `json:"success_submissions"`
}
