package requests

// SetCriteriaToGradeRequest request to set a criteria to a student's grade
type SetCriteriaToGradeRequest struct {
	CriteriaUUID  *string `json:"criteria_uuid" validate:"omitempty,uuid4"`
	ObjectiveUUID string  `json:"objective_uuid" validate:"required,uuid4"`
}

// GetSummarizedGradesInLaboratoryRequest request to get grade of an student in a laboratory
// that was graded with an specific rubric
type GetStudentGradeInLaboratoryWithRubricRequest struct {
	LaboratoryUUID string `validate:"required,uuid4"`
	RubricUUID     string `validate:"required,uuid4"`
	StudentUUID    string `validate:"required,uuid4"`
}

// SetCommentToGradeRequest request to set a comment to an student's grade
type SetCommentToGradeRequest struct {
	Comment string `json:"comment" validate:"required,min=8,max=510"`
}
