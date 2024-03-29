package dtos

type CreateRubricDTO struct {
	TeacherUUID string
	Name        string
}

type CreatedRubricDTO struct {
	UUID        string `json:"uuid"`
	TeacherUUID string `json:"-"`
	Name        string `json:"name"`
}

type DeleteRubricDTO struct {
	TeacherUUID string
	RubricUUID  string
}

type UpdateRubricNameDTO struct {
	TeacherUUID string
	RubricUUID  string
	Name        string
}
