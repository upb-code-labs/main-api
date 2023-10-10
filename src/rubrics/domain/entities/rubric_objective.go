package entities

type RubricObjective struct {
	UUID        string                    `json:"uuid"`
	RubricUUID  string                    `json:"-"`
	Description string                    `json:"description"`
	Criteria    []RubricObjectiveCriteria `json:"criteria"`
}
