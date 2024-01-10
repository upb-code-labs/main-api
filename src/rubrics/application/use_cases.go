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

func (useCases *RubricsUseCases) DeleteRubric(dto *dtos.DeleteRubricDTO) (err error) {
	// Check if the rubric belongs to the teacher
	teacherOwnsRubric, err := useCases.RubricsRepository.DoesTeacherOwnRubric(dto.TeacherUUID, dto.RubricUUID)
	if err != nil {
		return err
	}
	if !teacherOwnsRubric {
		return &errors.TeacherDoesNotOwnsRubric{}
	}

	// Delete the rubric
	err = useCases.RubricsRepository.Delete(dto.RubricUUID)
	if err != nil {
		return err
	}

	return nil
}

func (useCases *RubricsUseCases) UpdateRubricName(dto *dtos.UpdateRubricNameDTO) (err error) {
	// Check if the rubric belongs to the teacher
	teacherOwnsRubric, err := useCases.RubricsRepository.DoesTeacherOwnRubric(dto.TeacherUUID, dto.RubricUUID)
	if err != nil {
		return err
	}
	if !teacherOwnsRubric {
		return &errors.TeacherDoesNotOwnsRubric{}
	}

	// Update the name
	err = useCases.RubricsRepository.UpdateName(dto)
	if err != nil {
		return err
	}

	return nil
}

func (useCases *RubricsUseCases) AddObjectiveToRubric(dto *dtos.AddObjectiveToRubricDTO) (objectiveUUID string, err error) {
	// Check if the rubric belongs to the teacher
	teacherOwnsRubric, err := useCases.RubricsRepository.DoesTeacherOwnRubric(dto.TeacherUUID, dto.RubricUUID)
	if err != nil {
		return "", err
	}
	if !teacherOwnsRubric {
		return "", &errors.TeacherDoesNotOwnsRubric{}
	}

	// Add the objective
	objectiveUUID, err = useCases.RubricsRepository.AddObjectiveToRubric(dto.RubricUUID, dto.ObjectiveDescription)
	if err != nil {
		return "", err
	}

	return objectiveUUID, nil
}

func (useCases *RubricsUseCases) UpdateObjective(dto *dtos.UpdateObjectiveDTO) (err error) {
	// Check if the objective belongs to a rubric that belongs to the teacher
	teacherOwnsObjective, err := useCases.RubricsRepository.DoesTeacherOwnObjective(dto.TeacherUUID, dto.ObjectiveUUID)
	if err != nil {
		return err
	}
	if !teacherOwnsObjective {
		return &errors.TeacherDoesNotOwnsRubric{}
	}

	// Update the objective
	err = useCases.RubricsRepository.UpdateObjective(dto)
	if err != nil {
		return err
	}

	return nil
}

func (useCases *RubricsUseCases) DeleteObjective(dto *dtos.DeleteObjectiveDTO) (err error) {
	// Check if the objective belongs to a rubric that belongs to the teacher
	teacherOwnsObjective, err := useCases.RubricsRepository.DoesTeacherOwnObjective(dto.TeacherUUID, dto.ObjectiveUUID)
	if err != nil {
		return err
	}
	if !teacherOwnsObjective {
		return &errors.TeacherDoesNotOwnsRubric{}
	}

	// Delete the objective
	err = useCases.RubricsRepository.DeleteObjective(dto.ObjectiveUUID)
	if err != nil {
		return err
	}

	return nil
}

func (useCases *RubricsUseCases) AddCriteriaToObjective(dto *dtos.AddCriteriaToObjectiveDTO) (criteriaUUID string, err error) {
	// Check if the objective belongs to a rubric that belongs to the teacher
	teacherOwnsObjective, err := useCases.RubricsRepository.DoesTeacherOwnObjective(dto.TeacherUUID, dto.ObjectiveUUID)
	if err != nil {
		return "", err
	}
	if !teacherOwnsObjective {
		return "", &errors.TeacherDoesNotOwnsRubric{}
	}

	// Add the criteria
	criteriaUUID, err = useCases.RubricsRepository.AddCriteriaToObjective(dto)
	if err != nil {
		return "", err
	}

	return criteriaUUID, nil
}

func (useCases *RubricsUseCases) UpdateCriteria(dto *dtos.UpdateCriteriaDTO) (err error) {
	// Check if the criteria belongs to a rubric that belongs to the teacher
	teacherOwnsCriteria, err := useCases.RubricsRepository.DoesTeacherOwnCriteria(dto.TeacherUUID, dto.CriteriaUUID)
	if err != nil {
		return err
	}
	if !teacherOwnsCriteria {
		return &errors.TeacherDoesNotOwnsRubric{}
	}

	// Update the criteria
	err = useCases.RubricsRepository.UpdateCriteria(dto)
	if err != nil {
		return err
	}

	return nil
}

func (useCases *RubricsUseCases) DeleteCriteria(dto *dtos.DeleteCriteriaDTO) (err error) {
	// Check if the criteria belongs to a rubric that belongs to the teacher
	teacherOwnsCriteria, err := useCases.RubricsRepository.DoesTeacherOwnCriteria(dto.TeacherUUID, dto.CriteriaUUID)
	if err != nil {
		return err
	}
	if !teacherOwnsCriteria {
		return &errors.TeacherDoesNotOwnsRubric{}
	}

	// Delete the criteria
	err = useCases.RubricsRepository.DeleteCriteria(dto.CriteriaUUID)
	if err != nil {
		return err
	}

	return nil
}
