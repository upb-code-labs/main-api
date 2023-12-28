package implementations

import (
	"database/sql"

	"github.com/UPB-Code-Labs/main-api/src/languages/domain/entities"
	shared_infrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
)

type LanguagesRepository struct {
	Connection *sql.DB
}

// Singleton
var langRepositoryInstance *LanguagesRepository

func GetLanguagesRepositoryInstance() *LanguagesRepository {
	if langRepositoryInstance == nil {
		langRepositoryInstance = &LanguagesRepository{
			Connection: shared_infrastructure.GetPostgresConnection(),
		}
	}

	return langRepositoryInstance
}

// Methods implementation
func (repository *LanguagesRepository) GetAll() (languages []*entities.Language, err error) {
	return []*entities.Language{}, nil
}

func (repository *LanguagesRepository) GetByUUID(uuid string) (language *entities.Language, err error) {
	return &entities.Language{}, nil
}
