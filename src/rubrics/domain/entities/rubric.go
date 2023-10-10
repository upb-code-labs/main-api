package entities

type Rubric struct {
	UUID        string            `json:"uuid"`
	TeacherUUID string            `json:"-"`
	Name        string            `json:"name"`
	Objectives  []RubricObjective `json:"objectives"`
}
