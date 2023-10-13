package requests

type AddCriteriaToObjectiveRequest struct {
	Description string  `json:"description" validate:"required,min=8,max=510"`
	Weight      float64 `json:"weight" validate:"required,numeric,min=0,max=100"`
}
