package application

import (
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/dtos"
)

type LaboratoriesUseCases struct {
	Repository definitions.LaboratoriesRepository
}

func (useCases *LaboratoriesUseCases) CreateLaboratory(dto *dtos.CreateLaboratoryDTO) error {
	return useCases.Repository.SaveLaboratory(dto)
}
