package dtos

type AddCriteriaToObjectiveDTO struct {
	TeacherUUID         string
	ObjectiveUUID       string
	CriteriaDescription string
	CriteriaWeight      float64
}
