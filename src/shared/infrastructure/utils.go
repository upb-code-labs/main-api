package infrastructure

import (
	"encoding/json"
	"net/http"
	"time"

	sharedDomainErrors "github.com/UPB-Code-Labs/main-api/src/shared/domain/errors"
)

func ParseISODate(date string) (time.Time, error) {
	layout := "2006-01-02T15:04"
	return time.Parse(layout, date)
}

func ParseMicroserviceError(resp *http.Response, err error) error {
	if err != nil || resp.StatusCode != http.StatusOK {
		defaultErrorMessage := "There was an error while requesting the archives microservice"
		errorMessage := defaultErrorMessage

		// Decode the JSON from the body
		var responseJSON map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&responseJSON)
		if err != nil {
			return err
		}

		// Get the error message
		msg, ok := responseJSON["message"].(string)
		if ok {
			errorMessage = msg
		}

		// Return the error
		return &sharedDomainErrors.StaticFilesMicroserviceError{
			Code:    resp.StatusCode,
			Message: errorMessage,
		}
	}

	return nil
}
