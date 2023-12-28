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

func (useCases *LanguageUseCases) GetLanguageByUUID(uuid string) (*entities.Language, error) {
	return useCases.LanguageRepository.GetByUUID(uuid)
}
