package dtos

type AddObjectiveToRubricDTO struct {
	TeacherUUID          string
	RubricUUID           string
	ObjectiveDescription string
}

type UpdateObjectiveDTO struct {
	TeacherUUID        string
	ObjectiveUUID      string
	UpdatedDescription string
}

type DeleteObjectiveDTO struct {
	TeacherUUID   string `json:"teacher_uuid" binding:"required"`
	ObjectiveUUID string `json:"objective_uuid" binding:"required"`
}
