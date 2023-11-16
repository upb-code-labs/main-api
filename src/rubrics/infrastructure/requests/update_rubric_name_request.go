package requests

type UpdateRubricNameRequest struct {
	Name string `json:"name" validate:"required,min=4,max=96"`
}
