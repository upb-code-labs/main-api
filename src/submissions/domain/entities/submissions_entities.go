package entities

type Submission struct {
	UUID        string `json:"uuid"`
	ArchiveUUID string `json:"archive_uuid"`
	Passing     bool   `json:"passing"`
	Status      string `json:"status"`
	Stdout      string `json:"stdout"`
}

type SubmissionWork struct {
	SubmissionUUID        string `json:"submission_uuid"`
	LanguageUUID          string `json:"language_uuid"`
	SubmissionArchiveUUID string `json:"submission_archive_uuid"`
	TestArchiveUUID       string `json:"test_archive_uuid"`
}
