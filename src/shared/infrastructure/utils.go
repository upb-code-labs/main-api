package infrastructure

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	sharedDomainErrors "github.com/UPB-Code-Labs/main-api/src/shared/domain/errors"
	"github.com/gabriel-vasile/mimetype"
)

// ParseISODate parses a date in ISO format received from a date-time input
func ParseISODate(date string) (time.Time, error) {
	layout := "2006-01-02T15:04"
	return time.Parse(layout, date)
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

		// Decode the JSON from the body
		var responseJSON map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&responseJSON)
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
