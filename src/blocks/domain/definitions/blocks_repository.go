package definitions

import "mime/multipart"

type BlockRepository interface {
	UpdateMarkdownBlockContent(blockUUID string, content string) (err error)
	DoesTeacherOwnsMarkdownBlock(teacherUUID string, blockUUID string) (bool, error)

	SaveTestsArchive(file *multipart.File) (uuid string, err error)
}
