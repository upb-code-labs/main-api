package dtos

import "mime/multipart"

type CreateSubmissionDTO struct {
	StudentUUID       string
	TestBlockUUID     string
	SubmissionArchive *multipart.File
}

type GetSubmissionDTO struct {
	StudentUUID   string
	TestBlockUUID string
}
