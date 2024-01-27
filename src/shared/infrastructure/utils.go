package infrastructure

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"

	sharedDomainErrors "github.com/UPB-Code-Labs/main-api/src/shared/domain/errors"
	"github.com/gabriel-vasile/mimetype"
)

// ParseRFCEDate parses a date in RFC3339 format
func ParseRFCEDate(date string) (time.Time, error) {
	return time.Parse(time.RFC3339, date)
}

// ValidateMultipartFileHeader validates the multipart archive according to the
// environment configuration and domain rules
func ValidateMultipartFileHeader(multipartHeader *multipart.FileHeader) error {
	if multipartHeader.Size > GetEnvironment().ArchiveMaxSizeKb*1024 {
		return &sharedDomainErrors.GenericDomainError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("The archive must be less than %d KB", GetEnvironment().ArchiveMaxSizeKb),
		}
	}

	file, err := multipartHeader.Open()
	if err != nil {
		return &sharedDomainErrors.GenericDomainError{
			Code:    http.StatusInternalServerError,
			Message: "There was an error while reading the test archive",
		}
	}
	defer file.Close()

	mtype, err := mimetype.DetectReader(file)
	if err != nil {
		return &sharedDomainErrors.GenericDomainError{
			Code:    http.StatusInternalServerError,
			Message: "There was an error while reading the MIME type of the test archive",
		}
	}

	if mtype.String() != "application/zip" {
		return &sharedDomainErrors.GenericDomainError{
			Code:    http.StatusBadRequest,
			Message: "Please, make sure to send a ZIP archive",
		}
	}

	return nil
}

// ParseMicroserviceError parses the error returned by the archives microservice
func ParseMicroserviceError(resp *http.Response, err error) error {
	statusStr := fmt.Sprintf("%d", resp.StatusCode)
	isInTwoHundredsGroup := statusStr[0] == '2'

	if err != nil || !isInTwoHundredsGroup {
		defaultErrorMessage := "There was an error while requesting the archives microservice"
		errorMessage := defaultErrorMessage

		// Decode the body
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return &sharedDomainErrors.GenericDomainError{
				Code:    http.StatusBadRequest,
				Message: defaultErrorMessage,
			}
		}

		// Parse the JSON
		var responseJSON map[string]interface{}
		err = json.Unmarshal(body, &responseJSON)
		if err != nil {
			return &sharedDomainErrors.GenericDomainError{
				Code:    http.StatusBadRequest,
				Message: defaultErrorMessage,
			}
		}

		// Get the error message
		msg, ok := responseJSON["message"].(string)
		if ok {
			errorMessage = msg
		}

		// Return the error
		return &sharedDomainErrors.GenericDomainError{
			Code:    resp.StatusCode,
			Message: errorMessage,
		}
	}

	return nil
}

type baseMultipartFormBuffer struct {
	BodyBuffer       *bytes.Buffer
	BodyBufferWriter *multipart.Writer
}

func GetMultipartFormBufferFromFile(file *multipart.File) (*baseMultipartFormBuffer, error) {
	FILE_NAME := "archive.zip"
	FILE_CONTENT_TYPE := "application/zip"

	// Create multipart writer
	var bodyBuffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&bodyBuffer)

	// Add the file field to the request
	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition",
		fmt.Sprintf(
			`form-data; name="%s"; filename="%s"`,
			"file",
			FILE_NAME,
		),
	)
	header.Set("Content-Type", FILE_CONTENT_TYPE)

	// Reset the file pointer to the beginning
	_, err := (*file).Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Add the file to the request
	fileWriter, err := multipartWriter.CreatePart(header)
	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(fileWriter, *file); err != nil {
		return nil, err
	}

	return &baseMultipartFormBuffer{
		BodyBuffer:       &bodyBuffer,
		BodyBufferWriter: multipartWriter,
	}, nil
}
