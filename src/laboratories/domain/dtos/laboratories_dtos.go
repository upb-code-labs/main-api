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
	UserRole       string
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
	StudentsProgress []*SummarizedStudentProgressDTO `json:"students_progress"`
}

type SummarizedStudentProgressDTO struct {
	StudentUUID        string `json:"student_uuid"`
	StudentFullName    string `json:"student_full_name"`
	PendingSubmissions int    `json:"pending_submissions"`
	RunningSubmissions int    `json:"running_submissions"`
	FailingSubmissions int    `json:"failing_submissions"`
	SuccessSubmissions int    `json:"success_submissions"`
}

type GetProgressOfStudentInLaboratoryDTO struct {
	UserUUID       string
	UserRole       string
	LaboratoryUUID string
	StudentUUID    string
}

type StudentProgressInLaboratoryDTO struct {
	TotalTestBlocks    int                               `json:"total_test_blocks"`
	StudentSubmissions []*SummarizedStudentSubmissionDTO `json:"submissions"`
}

type SummarizedStudentSubmissionDTO struct {
	SubmissionUUID        string `json:"uuid"`
	SubmissionArchiveUUID string `json:"archive_uuid"`
	TestBlockName         string `json:"test_block_name"`
	SubmissionStatus      string `json:"status"`
	IsSubmissionPassing   bool   `json:"is_passing"`
}

type LaboratoryDetailsDTO struct {
	UUID        string  `json:"uuid"`
	CourseUUID  string  `json:"-"`
	RubricUUID  *string `json:"rubric_uuid"`
	Name        string  `json:"name"`
	OpeningDate string  `json:"opening_date"`
	DueDate     string  `json:"due_date"`
}
