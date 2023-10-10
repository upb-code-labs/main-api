package definitions

import (
	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/entities"
)

type RubricsRepository interface {
	Save(dto *dtos.CreateRubricDTO) (rubric *entities.Rubric, err error)
	GetByUUID(uuid string) (rubric *entities.Rubric, err error)
}