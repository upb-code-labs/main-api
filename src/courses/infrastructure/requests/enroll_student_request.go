package requests

type EnrollStudentRequest struct {
	StudentUUID string `json:"student_uuid" validate:"required,uuid4"`
}
