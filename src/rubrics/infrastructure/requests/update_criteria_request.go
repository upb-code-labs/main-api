package requests

type UpdateCriteriaRequest struct {
	Description string  `json:"description" validate:"required,min=8,max=510"`
	Weight      float64 `json:"weight" default:"0" validate:"numeric,min=0,max=100"`
}
