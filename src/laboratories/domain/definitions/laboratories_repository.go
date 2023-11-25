package definitions

import "github.com/UPB-Code-Labs/main-api/src/laboratories/domain/dtos"

type LaboratoriesRepository interface {
	SaveLaboratory(dto *dtos.CreateLaboratoryDTO) error
}
