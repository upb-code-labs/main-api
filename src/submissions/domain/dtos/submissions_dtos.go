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

type SubmissionStatusUpdateDTO struct {
	SubmissionUUID   string `json:"submission_uuid"`
	SubmissionStatus string `json:"submission_status"`
	TestsPassed      bool   `json:"tests_passed"`
	TestsOutput      string `json:"tests_output"`
}
