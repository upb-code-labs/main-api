package definitions

import (
	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/entities"
)

type RubricsRepository interface {
	Save(dto *dtos.CreateRubricDTO) (rubric *entities.Rubric, err error)
	GetByUUID(uuid string) (rubric *entities.Rubric, err error)
	GetAllCreatedByTeacher(teacherUUID string) (rubrics []*dtos.CreatedRubricDTO, err error)

	DoesTeacherOwnRubric(teacherUUID string, rubricUUID string) (bool, error)
	DoesTeacherOwnObjective(teacherUUID string, objectiveUUID string) (bool, error)
	DoesTeacherOwnCriteria(teacherUUID string, criteriaUUID string) (bool, error)

	UpdateName(dto *dtos.UpdateRubricNameDTO) (err error)

	AddObjectiveToRubric(rubricUUID string, objectiveDescription string) (objectiveUUID string, err error)
	UpdateObjective(dto *dtos.UpdateObjectiveDTO) (err error)

	AddCriteriaToObjective(dto *dtos.AddCriteriaToObjectiveDTO) (criteriaUUID string, err error)
	UpdateCriteria(dto *dtos.UpdateCriteriaDTO) (err error)
}
