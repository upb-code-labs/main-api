package application

import (
	"github.com/UPB-Code-Labs/main-api/src/languages/domain/definitions"
	"github.com/UPB-Code-Labs/main-api/src/languages/domain/entities"
)

type LanguageUseCases struct {
	LanguageRepository definitions.LanguagesRepository
}

func (useCases *LanguageUseCases) GetLanguages() ([]*entities.Language, error) {
	return useCases.LanguageRepository.GetAll()
}

func (useCases *LanguageUseCases) GetLanguageTemplate(uuid string) ([]byte, error) {
	// Get the information of the language from the database
	langTemplateUUID, err := useCases.LanguageRepository.GetTemplateArchiveUUIDByLanguageUUID(uuid)
	if err != nil {
		return nil, err
	}

	// Return an empty template bytes array
	return useCases.LanguageRepository.GetTemplateBytes(langTemplateUUID)
}
