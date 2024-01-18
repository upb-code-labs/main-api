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
	endpoint := fmt.Sprintf("/api/v1/submissions/test_blocks/%s", dto.blockUUID)

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

type closeNotifierRecorder struct {
	*httptest.ResponseRecorder
	closeNotify chan bool
}

func (c *closeNotifierRecorder) CloseNotify() <-chan bool {
	return c.closeNotify
}

type RealTimeSubmissionStatusResponse struct {
	w   *closeNotifierRecorder
	r   *http.Request
	err error
}

// GetRealTimeSubmissionStatus sends a request to the server to get the real time status of a submission and returns the response recorder
func GetRealTimeSubmissionStatus(testBlockUUID string, cookie *http.Cookie) *RealTimeSubmissionStatusResponse {
	// Create the request
	endpoint := fmt.Sprintf("/api/v1/submissions/test_blocks/%s/status", testBlockUUID)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return &RealTimeSubmissionStatusResponse{err: err}
	}

	req.AddCookie(cookie)

	// Send the request
	w := &closeNotifierRecorder{
		ResponseRecorder: httptest.NewRecorder(),
		closeNotify:      make(chan bool, 1),
	}
	router.ServeHTTP(w, req)

	return &RealTimeSubmissionStatusResponse{
		w:   w,
		r:   req,
		err: err,
	}
}

func GetSubmissionArchive(submissionUUID string, cookie *http.Cookie) (bytes []byte, statusCode int) {
	endpoint := fmt.Sprintf("/api/v1/submissions/%s/archive", submissionUUID)
	w, r := PrepareRequest("GET", endpoint, nil)
	r.AddCookie(cookie)
	router.ServeHTTP(w, r)
	return w.Body.Bytes(), w.Code
}
