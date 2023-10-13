package requests

type CreateRubricRequest struct {
	Name string `json:"name" validate:"required,min=4,max=96"`
}
