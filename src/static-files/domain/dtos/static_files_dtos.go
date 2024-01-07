package dtos

import "mime/multipart"

type SaveStaticFileDTO struct {
	FileType string `json:"archive_type"`
	File     *multipart.File
}

type OverwriteStaticFileDTO struct {
	FileUUID string `json:"archive_uuid"`
	FileType string `json:"archive_type"`
	File     *multipart.File
}

type StaticFileArchiveDTO struct {
	FileUUID string `json:"archive_uuid"`
	FileType string `json:"archive_type"`
}
