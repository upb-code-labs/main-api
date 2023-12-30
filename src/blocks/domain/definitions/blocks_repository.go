package definitions

import (
	"mime/multipart"

	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/dtos"
)

type BlockRepository interface {
	// Update the markdown text of a markdown block
	UpdateMarkdownBlockContent(blockUUID string, content string) (err error)

	// Functions to check blocks ownership
	DoesTeacherOwnsMarkdownBlock(teacherUUID string, blockUUID string) (bool, error)
	DoesTeacherOwnsTestBlock(teacherUUID string, blockUUID string) (bool, error)

	// Create a new test block
	SaveTestsArchive(file *multipart.File) (uuid string, err error)

	// Get the UUID of the `zip` archive saved in the static files microservice
	GetTestArchiveUUIDFromTestBlockUUID(blockUUID string) (uuid string, err error)

	// Overwrite the `zip` archive saved in the static files microservice
	OverwriteTestsArchive(uuid string, file *multipart.File) (err error)

	// Update the test block information in the database
	UpdateTestBlock(*dtos.UpdateTestBlockDTO) (err error)
}
