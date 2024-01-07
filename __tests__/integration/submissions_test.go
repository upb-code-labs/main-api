package integration

import (
	"bufio"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	submissionsDTOs "github.com/UPB-Code-Labs/main-api/src/submissions/domain/dtos"
	"github.com/stretchr/testify/require"
)

func TestSubmitSolutionToTestBlock(t *testing.T) {
	c := require.New(t)

	// ## Setup

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, status := CreateCourse("Submit solution to test block test - course")
	c.Equal(http.StatusCreated, status)

	// Create a laboratory
	laboratoryCreationResponse, status := CreateLaboratory(cookie, map[string]interface{}{
		"name":         "Submit solution to test block test - laboratory",
		"course_uuid":  courseUUID,
		"opening_date": "2023-12-01T08:00",
		"due_date":     "3023-12-01T00:00",
	})
	c.Equal(http.StatusCreated, status)
	laboratoryUUID := laboratoryCreationResponse["uuid"].(string)

	// Get the supported languages
	languagesResponse, status := GetSupportedLanguages(cookie)
	c.Equal(http.StatusOK, status)

	languages := languagesResponse["languages"].([]interface{})
	c.Greater(len(languages), 0)

	firstLanguage := languages[0].(map[string]interface{})
	firstLanguageUUID := firstLanguage["uuid"].(string)

	// Create a test block
	zipFile, err := GetSampleTestsArchive()
	c.Nil(err)

	blockCreationResponse, status := CreateTestBlock(&CreateTestBlockUtilsDTO{
		laboratoryUUID: laboratoryUUID,
		languageUUID:   firstLanguageUUID,
		blockName:      "Update test block test - block",
		cookie:         cookie,
		testFile:       zipFile,
	})
	c.Equal(http.StatusCreated, status)
	testBlockUUID := blockCreationResponse["uuid"].(string)

	// Add a student to the course
	invitationCode, code := GetInvitationCode(courseUUID)
	c.Equal(http.StatusOK, code)
	c.NotEmpty(invitationCode)

	_, code = AddStudentToCourse(invitationCode)
	c.Equal(http.StatusOK, code)

	// ## Test
	zipFile, err = GetSampleSubmissionArchive()
	c.Nil(err)

	// Login as a student
	w, r = PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredStudentEmail,
		"password": registeredStudentPass,
	})
	router.ServeHTTP(w, r)
	cookie = w.Result().Cookies()[0]

	// Submit the solution
	submissionResponse, status := SubmitSolutionToTestBlock(&SubmitSToTestBlockUtilsDTO{
		blockUUID: testBlockUUID,
		cookie:    cookie,
		file:      zipFile,
	})

	c.Equal(http.StatusCreated, status)
	c.NotEmpty(submissionResponse["uuid"])

	// ## Get submission status
	submissionUUID := submissionResponse["uuid"].(string)

	response := GetRealTimeSubmissionStatus(testBlockUUID, cookie)
	c.Nil(response.err)

	// Receive events
	EXPECTED_EVENTS_COUNT := 3
	receivedEventsCount := 0
	receivedEvents := make(
		[]*submissionsDTOs.SubmissionStatusUpdateDTO,
		EXPECTED_EVENTS_COUNT,
	)

	scanner := bufio.NewScanner(response.w.Body)
	for scanner.Scan() {
		// Read the rew text
		response := scanner.Text()

		// Check if it is a data event
		isData := strings.HasPrefix(response, "data:")
		if !isData {
			continue
		}

		// Remove the prefix
		response = strings.TrimPrefix(response, "data:")
		response = strings.TrimSpace(response)

		var event submissionsDTOs.SubmissionStatusUpdateDTO
		err := json.Unmarshal([]byte(response), &event)
		c.Nil(err)

		// Add the event to the list
		receivedEvents[receivedEventsCount] = &event

		receivedEventsCount++
		if receivedEventsCount == EXPECTED_EVENTS_COUNT {
			break
		}
	}

	c.Nil(scanner.Err())

	// Check the events
	partialEventsStatus := []string{
		"pending",
		"running",
	}

	c.Equal(EXPECTED_EVENTS_COUNT, receivedEventsCount)
	for idx, event := range receivedEvents {
		c.Equal(submissionUUID, event.SubmissionUUID)

		if idx < 2 {
			c.Equal(partialEventsStatus[idx], event.SubmissionStatus)
			c.False(event.TestsPassed)
			c.Empty(event.TestsOutput)
		} else {
			c.Equal("ready", event.SubmissionStatus)
			c.True(event.TestsPassed)
			c.NotEmpty(event.TestsOutput)
		}
	}
}
