package dtos

type GetCourseLaboratoriesDTO struct {
	CourseUUID string
	UserUUID   string
	UserRole   string
}

type BaseLaboratoryDTO struct {
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	OpeningDate string `json:"opening_date"`
	DueDate     string `json:"due_date"`
}
