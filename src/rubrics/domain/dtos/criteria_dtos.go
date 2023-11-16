package dtos

type AddCriteriaToObjectiveDTO struct {
	TeacherUUID         string
	ObjectiveUUID       string
	CriteriaDescription string
	CriteriaWeight      float64
}

type UpdateCriteriaDTO struct {
	TeacherUUID         string
	CriteriaUUID        string
	CriteriaDescription string
	CriteriaWeight      float64
}
