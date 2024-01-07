package definitions

import (
	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/dtos"
)

type BlockRepository interface {
	// Update the markdown text of a markdown block
	UpdateMarkdownBlockContent(blockUUID string, content string) (err error)

	// Check blocks ownership
	DoesTeacherOwnsMarkdownBlock(teacherUUID string, blockUUID string) (bool, error)
	DoesTeacherOwnsTestBlock(teacherUUID string, blockUUID string) (bool, error)

	// Check blocks permissions
	CanStudentSubmitToTestBlock(studentUUID string, testBlockUUID string) (bool, error)

	// Get the UUID of the `zip` archive saved in the static files microservice
	GetTestArchiveUUIDFromTestBlockUUID(blockUUID string) (uuid string, err error)

	// Update the test block information in the database
	UpdateTestBlock(*dtos.UpdateTestBlockDTO) (err error)

	// Delete blocks
	DeleteMarkdownBlock(blockUUID string) (err error)
	DeleteTestBlock(blockUUID string) (err error)

	// Get the laboratory the block belongs to
	GetTestBlockLaboratoryUUID(blockUUID string) (laboratoryUUID string, err error)
}
