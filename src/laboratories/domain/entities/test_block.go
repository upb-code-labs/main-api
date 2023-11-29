package entities

type TestBlock struct {
	UUID            string `json:"uuid"`
	LanguageUUID    string `json:"language_uuid"`
	TestArchiveUUID string `json:"test_archive_uuid"`
	Name            string `json:"name"`
	Index           int    `json:"index"`
}
