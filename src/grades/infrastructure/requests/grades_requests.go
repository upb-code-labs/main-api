package requests

// SetCriteriaToGradeRequest request to set a criteria to a student's grade
type SetCriteriaToGradeRequest struct {
	CriteriaUUID  *string `json:"criteria_uuid" validate:"omitempty,uuid4"`
	ObjectiveUUID string  `json:"objective_uuid" validate:"required,uuid4"`
	RubricUUID    string  `json:"rubric_uuid" validate:"required,uuid4"`
}

type SetCriteriaToGradeRequestUUIDs struct {
	StudentUUID    string `validate:"required,uuid4"`
	LaboratoryUUID string `validate:"required,uuid4"`
}
