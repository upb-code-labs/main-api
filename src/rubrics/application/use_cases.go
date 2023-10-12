package application

import (
	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/entities"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/errors"
)

type RubricsUseCases struct {
	RubricsRepository definitions.RubricsRepository
}

func (useCases *RubricsUseCases) CreateRubric(dto *dtos.CreateRubricDTO) (rubric *entities.Rubric, err error) {
	savedRubric, err := useCases.RubricsRepository.Save(dto)
	if err != nil {
		return nil, err
	}

	return savedRubric, nil
}

func (useCases *RubricsUseCases) GetRubricsCreatedByTeacher(teacherUUID string) (rubrics []*dtos.CreatedRubricDTO, err error) {
	rubrics, err = useCases.RubricsRepository.GetAllCreatedByTeacher(teacherUUID)
	if err != nil {
		return nil, err
	}

	return rubrics, nil
}

func (useCases *RubricsUseCases) GetRubricByUUID(dto *dtos.GetRubricDto) (rubric *entities.Rubric, err error) {
	// Get the rubric
	rubric, err = useCases.RubricsRepository.GetByUUID(dto.RubricUUID)
	if err != nil {
		return nil, err
	}

	// Check if the rubric belongs to the teacher
	if rubric.TeacherUUID != dto.TeacherUUID {
		return nil, &errors.TeacherDoesNotOwnsRubric{}
	}

	return rubric, nil
}

func (useCases *RubricsUseCases) AddObjectiveToRubric(dto *dtos.AddObjectiveToRubricDTO) (objectiveUUID string, err error) {
	// Get the rubric
	rubric, err := useCases.RubricsRepository.GetByUUID(dto.RubricUUID)
	if err != nil {
		return "", err
	}

	// Check if the rubric belongs to the teacher
	if rubric.TeacherUUID != dto.TeacherUUID {
		return "", &errors.TeacherDoesNotOwnsRubric{}
	}

	// Add the objective
	objectiveUUID, err = useCases.RubricsRepository.AddObjectiveToRubric(dto.RubricUUID, dto.ObjectiveDescription)
	if err != nil {
		return "", err
	}

	return objectiveUUID, nil
}
