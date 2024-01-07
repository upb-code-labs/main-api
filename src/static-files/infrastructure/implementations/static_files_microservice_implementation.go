package implementations

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	sharedDomainErrors "github.com/UPB-Code-Labs/main-api/src/shared/domain/errors"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/UPB-Code-Labs/main-api/src/static-files/domain/dtos"
)

type StaticFilesMicroserviceImplementation struct{}

// SaveArchive saves a file in the static files microservice
func (implementation *StaticFilesMicroserviceImplementation) SaveArchive(dto *dtos.SaveStaticFileDTO) (fileUUID string, err error) {
	// Create multipart writer
	staticFilesEndpoint := fmt.Sprintf(
		"%s/archives/save",
		sharedInfrastructure.GetEnvironment().StaticFilesMicroserviceAddress,
	)

	baseMultipartBuffer, err := sharedInfrastructure.GetMultipartFormBufferFromFile(dto.File)
	if err != nil {
		return "", err
	}

	bodyBufferWriter := baseMultipartBuffer.BodyBufferWriter
	bodyBuffer := baseMultipartBuffer.BodyBuffer

	// Add the file type field to the request
	err = bodyBufferWriter.WriteField("archive_type", dto.FileType)
	if err != nil {
		return "", err
	}

	// Close the writer
	err = bodyBufferWriter.Close()
	if err != nil {
		return "", err
	}

	// Prepare the request
	req, err := http.NewRequest("POST", staticFilesEndpoint, bodyBuffer)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", bodyBufferWriter.FormDataContentType())

	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)

	// Forward error message if any
	microserviceError := sharedInfrastructure.ParseMicroserviceError(res, err)
	if microserviceError != nil {
		return "", microserviceError
	}

	// Return the UUID of the saved file
	var response map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	if response["uuid"] == nil {
		return "", &sharedDomainErrors.GenericDomainError{
			Code:    http.StatusInternalServerError,
			Message: "The static files microservice did not return the UUID of the saved file",
		}
	}

	return response["uuid"].(string), nil
}

// OverwriteArchive overwrites a file in the static files microservice
func (implementation *StaticFilesMicroserviceImplementation) OverwriteArchive(dto *dtos.OverwriteStaticFileDTO) error {
	// Create multipart writer
	staticFilesEndpoint := fmt.Sprintf(
		"%s/archives/overwrite",
		sharedInfrastructure.GetEnvironment().StaticFilesMicroserviceAddress,
	)

	baseMultipartBuffer, err := sharedInfrastructure.GetMultipartFormBufferFromFile(dto.File)
	if err != nil {
		return err
	}

	// Add the file type field to the request
	err = baseMultipartBuffer.BodyBufferWriter.WriteField("archive_type", dto.FileType)
	if err != nil {
		return err
	}

	// Add the archive uuid field to the request
	err = baseMultipartBuffer.BodyBufferWriter.WriteField("archive_uuid", dto.FileUUID)
	if err != nil {
		return err
	}

	// Close the writer
	err = baseMultipartBuffer.BodyBufferWriter.Close()
	if err != nil {
		return err
	}

	// Prepare the request
	req, err := http.NewRequest("PUT", staticFilesEndpoint, baseMultipartBuffer.BodyBuffer)

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", baseMultipartBuffer.BodyBufferWriter.FormDataContentType())

	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)

	// Forward error message if any
	microserviceError := sharedInfrastructure.ParseMicroserviceError(res, err)
	if microserviceError != nil {
		return microserviceError
	}

	return nil
}

// DeleteArchive deletes a file in the static files microservice
func (implementation *StaticFilesMicroserviceImplementation) GetArchiveBytes(dto *dtos.StaticFileArchiveDTO) ([]byte, error) {
	// Prepare the request
	staticFilesEndpoint := fmt.Sprintf(
		"%s/archives/download",
		sharedInfrastructure.GetEnvironment().StaticFilesMicroserviceAddress,
	)

	// Create request payload from the dto
	body, err := json.Marshal(dto)
	if err != nil {
		return []byte{}, err
	}

	// Create request
	request, err := http.NewRequest(
		"POST",
		staticFilesEndpoint,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return []byte{}, err
	}

	request.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)

	// Forward error message if any
	microserviceError := sharedInfrastructure.ParseMicroserviceError(response, err)
	if microserviceError != nil {
		return []byte{}, microserviceError
	}

	defer response.Body.Close()

	// Handle error
	if response.StatusCode != http.StatusOK {
		return []byte{}, errors.New(
			"there was an error while trying to get the archive from the static files microservice",
		)
	}

	// Read response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	// Return response
	return responseBody, nil
}

// GetLanguageTemplateArchiveBytes gets a language template file from the static files microservice
func (implementation *StaticFilesMicroserviceImplementation) GetLanguageTemplateArchiveBytes(languageUUID string) ([]byte, error) {
	// Prepare the request
	staticFilesEndpoint := fmt.Sprintf(
		"%s/templates/%s",
		sharedInfrastructure.GetEnvironment().StaticFilesMicroserviceAddress,
		languageUUID,
	)

	// Create request
	request, err := http.NewRequest(
		"GET",
		staticFilesEndpoint,
		nil,
	)
	if err != nil {
		return []byte{}, err
	}

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)

	// Forward error message if any
	microserviceError := sharedInfrastructure.ParseMicroserviceError(response, err)
	if microserviceError != nil {
		return []byte{}, microserviceError
	}

	defer response.Body.Close()

	// Handle error
	if response.StatusCode != http.StatusOK {
		return []byte{}, errors.New(
			"there was an error while trying to get the archive from the static files microservice",
		)
	}

	// Read response
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return []byte{}, err
	}

	// Return response
	return responseBody, nil
}

// DeleteArchive deletes a file in the static files microservice
func (implementation *StaticFilesMicroserviceImplementation) DeleteArchive(dto *dtos.StaticFileArchiveDTO) error {
	// Prepare the request
	staticFilesEndpoint := fmt.Sprintf(
		"%s/archives/delete",
		sharedInfrastructure.GetEnvironment().StaticFilesMicroserviceAddress,
	)

	// Create request payload from the dto
	body, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	// Create request
	request, err := http.NewRequest(
		"POST",
		staticFilesEndpoint,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	response, err := client.Do(request)

	// Forward error message if any
	microserviceError := sharedInfrastructure.ParseMicroserviceError(response, err)
	if microserviceError != nil {
		return microserviceError
	}

	defer response.Body.Close()

	// Handle error
	if response.StatusCode != http.StatusOK {
		return errors.New(
			"there was an error while trying to delete the archive from the static files microservice",
		)
	}

	// Return response
	return nil
}
