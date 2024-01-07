package definitions

import "github.com/UPB-Code-Labs/main-api/src/static-files/domain/dtos"

type StaticFilesRepository interface {
	SaveArchive(dto *dtos.SaveStaticFileDTO) (fileUUID string, err error)
	OverwriteArchive(dto *dtos.OverwriteStaticFileDTO) error

	GetArchiveBytes(dto *dtos.StaticFileArchiveDTO) ([]byte, error)
	GetLanguageTemplateArchiveBytes(languageUUID string) ([]byte, error)

	DeleteArchive(dto *dtos.StaticFileArchiveDTO) error
}
