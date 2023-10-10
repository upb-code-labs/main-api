package application

import (
	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/entities"
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
