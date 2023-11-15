package requests

type UpdateObjectiveRequest struct {
	Description string `json:"description" validate:"required,min=8,max=510"`
}
