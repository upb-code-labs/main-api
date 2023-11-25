package implementations

import (
	"database/sql"

	"github.com/UPB-Code-Labs/main-api/src/laboratories/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
)

type LaboratoriesPostgresRepository struct {
	Connection *sql.DB
}

// Singleton
var laboratoriesPostgresRepositoryInstance *LaboratoriesPostgresRepository

func GetLaboratoriesPostgresRepositoryInstance() *LaboratoriesPostgresRepository {
	if laboratoriesPostgresRepositoryInstance == nil {
		laboratoriesPostgresRepositoryInstance = &LaboratoriesPostgresRepository{
			Connection: infrastructure.GetPostgresConnection(),
		}
	}

	return laboratoriesPostgresRepositoryInstance
}

func (repository *LaboratoriesPostgresRepository) SaveLaboratory(dto *dtos.CreateLaboratoryDTO) error {
	return nil
}
