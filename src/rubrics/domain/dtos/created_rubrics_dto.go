package dtos

type CreatedRubricDTO struct {
	UUID        string `json:"uuid"`
	TeacherUUID string `json:"-"`
	Name        string `json:"name"`
}
