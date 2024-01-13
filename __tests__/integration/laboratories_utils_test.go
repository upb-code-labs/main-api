package integration

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
)

func CreateLaboratory(cookie *http.Cookie, payload map[string]interface{}) (response map[string]interface{}, statusCode int) {
	w, r := PrepareRequest("POST", "/api/v1/laboratories", payload)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func GetLaboratoryByUUID(cookie *http.Cookie, uuid string) (response map[string]interface{}, statusCode int) {
	w, r := PrepareRequest("GET", "/api/v1/laboratories/"+uuid, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func GetLaboratoryInformationByUUID(cookie *http.Cookie, uuid string) (response map[string]interface{}, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/laboratories/%s/information", uuid)
	w, r := PrepareRequest("GET", endpoint, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func UpdateLaboratory(cookie *http.Cookie, uuid string, payload map[string]interface{}) (response map[string]interface{}, statusCode int) {
	w, r := PrepareRequest("PUT", "/api/v1/laboratories/"+uuid, payload)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func CreateMarkdownBlock(cookie *http.Cookie, laboratoryUUID string) (response map[string]interface{}, statusCode int) {
	w, r := PrepareRequest("POST", "/api/v1/laboratories/markdown_blocks/"+laboratoryUUID, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

type CreateTestBlockUtilsDTO struct {
	laboratoryUUID string
	languageUUID   string
	blockName      string
	cookie         *http.Cookie
	testFile       *os.File
}

func CreateTestBlock(dto *CreateTestBlockUtilsDTO) (response map[string]interface{}, statusCode int) {
	// Create the request body
	var body bytes.Buffer

	// Create the multipart form
	writer := multipart.NewWriter(&body)

	// Add the file to the form
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", "form-data; name=\"test_archive\"; filename=\"test.zip\"")
	h.Set("Content-Type", "application/zip")

	fileWriter, err := writer.CreatePart(h)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(fileWriter, dto.testFile)
	if err != nil {
		panic(err)
	}

	// Add the text fields to the form
	err = writer.WriteField("block_name", dto.blockName)
	if err != nil {
		panic(err)
	}

	err = writer.WriteField("language_uuid", dto.languageUUID)
	if err != nil {
		panic(err)
	}

	// Close the multipart form
	err = writer.Close()
	if err != nil {
		panic(err)
	}

	// Create the request
	w, r := PrepareMultipartRequest("POST", "/api/v1/laboratories/test_blocks/"+dto.laboratoryUUID, &body)
	r.AddCookie(dto.cookie)
	r.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request
	router.ServeHTTP(w, r)
	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func GetStudentsProgressInLaboratory(laboratoryUUID string, cookie *http.Cookie) (response map[string]interface{}, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/laboratories/%s/progress", laboratoryUUID)
	w, r := PrepareRequest("GET", endpoint, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}
