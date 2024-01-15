package dtos

// GetSummarizedGradesInLaboratoryDTO data transfer object to parse the request of the endpoint
type GetSummarizedGradesInLaboratoryDTO struct {
	TeacherUUID    string
	LaboratoryUUID string
}

// SummarizedStudentGradeDTO data transfer object to be used as the response of the endpoint
type SummarizedStudentGradeDTO struct {
	StudentUUID     string  `json:"student_uuid"`
	StudentFullName string  `json:"student_full_name"`
	Grade           float64 `json:"grade"`
}
