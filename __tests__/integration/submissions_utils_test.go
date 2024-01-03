package integration

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
)

type SubmitSToTestBlockUtilsDTO struct {
	blockUUID string
	cookie    *http.Cookie
	file      *os.File
}

func SubmitSolutionToTestBlock(dto *SubmitSToTestBlockUtilsDTO) (response map[string]interface{}, statusCode int) {
	// Create the request body
	var body bytes.Buffer

	// Create the multipart form
	writer := multipart.NewWriter(&body)

	// Add the test file
	part, err := writer.CreateFormFile("submission_archive", dto.file.Name())
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(part, dto.file)
	if err != nil {
		panic(err)
	}

	// Close the multipart form
	err = writer.Close()
	if err != nil {
		panic(err)
	}

	// Create the request
	endpoint := fmt.Sprintf("/api/v1/submissions/%s", dto.blockUUID)

	req, err := http.NewRequest("POST", endpoint, &body)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", writer.FormDataContentType())
	req.AddCookie(dto.cookie)

	// Send the request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Parse the response
	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}
