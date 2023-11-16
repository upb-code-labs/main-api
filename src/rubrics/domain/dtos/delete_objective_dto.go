package dtos

type DeleteObjectiveDTO struct {
	TeacherUUID   string `json:"teacher_uuid" binding:"required"`
	ObjectiveUUID string `json:"objective_uuid" binding:"required"`
}
