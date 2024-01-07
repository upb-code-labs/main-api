package definitions

import "github.com/UPB-Code-Labs/main-api/src/languages/domain/entities"

type LanguagesRepository interface {
	GetAll() (languages []*entities.Language, err error)
	GetByUUID(uuid string) (language *entities.Language, err error)
	GetTemplateArchiveUUIDByLanguageUUID(uuid string) (templateUUID string, err error)
}
