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

func UpdateMarkdownBlockContent(cookie *http.Cookie, blockUUID string, payload map[string]interface{}) (response map[string]interface{}, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/blocks/markdown_blocks/%s/content", blockUUID)
	w, r := PrepareRequest("PATCH", endpoint, payload)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func DeleteMarkdownBlock(cookie *http.Cookie, blockUUID string) (response map[string]interface{}, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/blocks/markdown_blocks/%s", blockUUID)
	w, r := PrepareRequest("DELETE", endpoint, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

type UpdateTestBlockUtilsDTO struct {
	blockUUID    string
	languageUUID string
	blockName    string
	cookie       *http.Cookie
	testFile     *os.File
}

func UpdateTestBlock(dto *UpdateTestBlockUtilsDTO) (response map[string]interface{}, statusCode int) {
	// Create the request body
	var body bytes.Buffer

	// Create the multipart form
	writer := multipart.NewWriter(&body)

	// Add the block name
	_ = writer.WriteField("block_name", dto.blockName)

	// Add the language UUID
	_ = writer.WriteField("language_uuid", dto.languageUUID)

	// Add the test file
	if dto.testFile != nil {
		part, err := writer.CreateFormFile("test_archive", dto.testFile.Name())
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(part, dto.testFile)
		if err != nil {
			panic(err)
		}
	}

	// Close the multipart form
	err := writer.Close()
	if err != nil {
		panic(err)
	}

	// Create the request
	endpoint := fmt.Sprintf("/api/v1/blocks/test_blocks/%s", dto.blockUUID)
	req, err := http.NewRequest("PUT", endpoint, &body)
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

func DeleteTestBlock(cookie *http.Cookie, blockUUID string) (response map[string]interface{}, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/blocks/test_blocks/%s", blockUUID)
	w, r := PrepareRequest("DELETE", endpoint, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

type SwapBlocksUtilsDTO struct {
	FirstBlockUUID  string
	SecondBlockUUID string
	Cookie          *http.Cookie
}

func SwapBlocks(dto *SwapBlocksUtilsDTO) (response map[string]interface{}, statusCode int) {
	payload := map[string]interface{}{
		"first_block_uuid":  dto.FirstBlockUUID,
		"second_block_uuid": dto.SecondBlockUUID,
	}

	w, r := PrepareRequest("PATCH", "/api/v1/blocks/swap_index", payload)
	r.AddCookie(dto.Cookie)
	router.ServeHTTP(w, r)

	jsonResponse := ParseJsonResponse(w.Body)
	return jsonResponse, w.Code
}

func GetTestsArchive(testBlockUUID string, cookie *http.Cookie) (bytes []byte, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/blocks/test_blocks/%s/tests_archive", testBlockUUID)
	w, r := PrepareRequest("GET", endpoint, nil)
	r.AddCookie(cookie)

	router.ServeHTTP(w, r)
	return w.Body.Bytes(), w.Code
}
