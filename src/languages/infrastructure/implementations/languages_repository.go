package implementations

import (
	"context"
	"database/sql"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT 
		id, template_archive_id, name 
		FROM languages
	`

	rows, err := repository.Connection.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Parse the rows
	for rows.Next() {
		var language entities.Language
		err := rows.Scan(&language.UUID, &language.TemplateArchiveUUID, &language.Name)
		if err != nil {
			return nil, err
		}

		languages = append(languages, &language)
	}

	return languages, nil
}

func (repository *LanguagesRepository) GetByUUID(uuid string) (language *entities.Language, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT 
		id, template_archive_uuid, name 
		FROM languages
		WHERE uuid = $1
	`
	row := repository.Connection.QueryRowContext(ctx, query, uuid)

	// Parse the row
	err = row.Scan(&language.UUID, &language.TemplateArchiveUUID, &language.Name)
	if err != nil {
		return nil, err
	}

	return language, nil
}
