package dtos

import "mime/multipart"

type UpdateMarkdownBlockContentDTO struct {
	TeacherUUID string
	BlockUUID   string
	Content     string
}

type UpdateTestBlockDTO struct {
	TeacherUUID    string
	BlockUUID      string
	LanguageUUID   string
	Name           string
	NewTestArchive *multipart.File
}
