package requests

type CreateCourseRequest struct {
	Name string `json:"name" validate:"required,min=4,max=96"`
}

type EnrollStudentRequest struct {
	StudentUUID string `json:"student_uuid" validate:"required,uuid4"`
}
