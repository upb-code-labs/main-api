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

type DeleteBlockDTO struct {
	TeacherUUID string
	BlockUUID   string
}

type SwapBlocksDTO struct {
	TeacherUUID     string
	FirstBlockUUID  string
	SecondBlockUUID string
}
