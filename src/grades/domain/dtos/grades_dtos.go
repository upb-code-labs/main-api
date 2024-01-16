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

// SetCriteriaToGradeDTO data transfer object to parse the request of the endpoint
type SetCriteriaToGradeDTO struct {
	TeacherUUID    string
	LaboratoryUUID string
	StudentUUID    string
	RubricUUID     string
	CriteriaUUID   *string
	ObjectiveUUID  string
}

// CheckIfStudentHasGradeDTO data transfer object to check if a student has a grade in a laboratory
type CheckIfStudentHasGradeDTO struct {
	StudentUUID    string
	LaboratoryUUID string
	RubricUUID     string
}

type CreateStudentGradeDTO struct {
	CheckIfStudentHasGradeDTO
}

type GetStudentGradeDTO struct {
	CheckIfStudentHasGradeDTO
}

// GetStudentGradeInLaboratoryWithRubricDTO data transfer object to parse the request of the endpoint
type GetStudentGradeInLaboratoryWithRubricDTO struct {
	UserUUID       string
	StudentUUID    string
	LaboratoryUUID string
	RubricUUID     string
}

// STudentGradeInLaboratoryWithRubricDTO data transfer object to be used as the response of the endpoint
type StudentGradeInLaboratoryWithRubricDTO struct {
	Grade            float64                              `json:"grade"`
	Comment          string                               `json:"comment"`
	SelectedCriteria []*SelectedCriteriaInStudentGradeDTO `json:"selected_criteria"`
}

// SelectedCriteriaInStudentGradeDTO data transfer object to obtain the selected criteria in a student's grade
type SelectedCriteriaInStudentGradeDTO struct {
	ObjectiveUUID string `json:"objective_uuid"`
	CriteriaUUID  string `json:"criteria_uuid"`
}
