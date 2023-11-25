package definitions

import (
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/entities"
)

type LaboratoriesRepository interface {
	GetLaboratoryByUUID(uuid string) (laboratory *entities.Laboratory, err error)
	SaveLaboratory(dto *dtos.CreateLaboratoryDTO) (laboratory *entities.Laboratory, err error)
}