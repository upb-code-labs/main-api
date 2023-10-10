package entities

type RubricObjectiveCriteria struct {
	UUID          string  `json:"uuid"`
	ObjectiveUUID string  `json:"-"`
	Description   string  `json:"description"`
	Weight        float32 `json:"weight"`
}
