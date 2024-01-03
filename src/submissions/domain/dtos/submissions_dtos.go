package dtos

import "mime/multipart"

type CreateSubmissionDTO struct {
	StudentUUID       string
	TestBlockUUID     string
	SubmissionArchive *multipart.File
	SavedArchiveUUID  string
}

type GetSubmissionDTO struct {
	StudentUUID   string
	TestBlockUUID string
}
