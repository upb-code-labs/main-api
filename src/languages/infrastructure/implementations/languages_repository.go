package implementations

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/UPB-Code-Labs/main-api/src/languages/domain/entities"
	"github.com/UPB-Code-Labs/main-api/src/languages/domain/errors"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
)

type LanguagesRepository struct {
	Connection *sql.DB
}

// Singleton
var langRepositoryInstance *LanguagesRepository

func GetLanguagesRepositoryInstance() *LanguagesRepository {
	if langRepositoryInstance == nil {
		langRepositoryInstance = &LanguagesRepository{
			Connection: sharedInfrastructure.GetPostgresConnection(),
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
		id, template_archive_id, name 
		FROM languages
		WHERE id = $1
	`

	row := repository.Connection.QueryRowContext(ctx, query, uuid)

	// Parse the row
	language = &entities.Language{}
	err = row.Scan(&language.UUID, &language.TemplateArchiveUUID, &language.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &errors.LangNotFoundError{}
		}

		return nil, err
	}

	return language, nil
}

func (repository *LanguagesRepository) GetTemplateArchiveUUIDByLanguageUUID(uuid string) (templateUUID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `
		SELECT file_id
		FROM archives
		WHERE id = (
			SELECT
			template_archive_id
			FROM languages
			WHERE id = $1
		)
	`

	row := repository.Connection.QueryRowContext(ctx, query, uuid)

	// Parse the row
	err = row.Scan(&templateUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", &errors.LangNotFoundError{}
		}

		return "", err
	}

	return templateUUID, nil
}

func (repository *LanguagesRepository) GetTemplateBytes(uuid string) (template []byte, err error) {
	// Send a request to the static files microservice
	staticFilesMsEndpoint := fmt.Sprintf("%s/templates/%s", sharedInfrastructure.GetEnvironment().StaticFilesMicroserviceAddress, uuid)
	resp, err := http.Get(staticFilesMsEndpoint)

	// If there is an error try to forward the error message
	microserviceError := sharedInfrastructure.ParseMicroserviceError(resp, err)
	if microserviceError != nil {
		return nil, microserviceError
	}

	// Read the body
	defer resp.Body.Close()
	template, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return template, nil
}
