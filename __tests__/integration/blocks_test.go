package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateMarkdownBlockContent(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, status := CreateCourse("Update markdown block content test - course")
	c.Equal(http.StatusCreated, status)

	// Create a laboratory
	laboratoryCreationResponse, status := CreateLaboratory(cookie, map[string]interface{}{
		"name":         "Update markdown block content test - laboratory",
		"course_uuid":  courseUUID,
		"opening_date": defaultLaboratoryOpeningDate,
		"due_date":     defaultLaboratoryDueDate,
	})
	c.Equal(http.StatusCreated, status)
	laboratoryUUID := laboratoryCreationResponse["uuid"].(string)

	// Create a markdown block
	blockCreationResponse, status := CreateMarkdownBlock(cookie, laboratoryUUID)
	c.Equal(http.StatusCreated, status)
	markdownBlockUUID := blockCreationResponse["uuid"].(string)

	// Define tests cases
	testCases := []GenericTestCase{
		{
			Payload: map[string]interface{}{
				"content": "",
			},
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Payload: map[string]interface{}{
				"content": "# Updated main title",
			},
			ExpectedStatusCode: http.StatusNoContent,
		},
	}

	// Run tests
	for _, tc := range testCases {
		_, status := UpdateMarkdownBlockContent(cookie, markdownBlockUUID, tc.Payload)
		c.Equal(tc.ExpectedStatusCode, status)
	}

	// Verify that the content was updated
	response, status := GetLaboratoryByUUID(cookie, laboratoryUUID)
	c.Equal(http.StatusOK, status)
	c.Equal(1, len(response["markdown_blocks"].([]interface{})))

	block := response["markdown_blocks"].([]interface{})[0].(map[string]interface{})
	c.Equal("# Updated main title", block["content"].(string))
}

func TestDeleteMarkdownBlock(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, status := CreateCourse("Delete markdown block test - course")
	c.Equal(http.StatusCreated, status)

	// Create a laboratory
	laboratoryCreationResponse, status := CreateLaboratory(cookie, map[string]interface{}{
		"name":         "Delete markdown block test - laboratory",
		"course_uuid":  courseUUID,
		"opening_date": defaultLaboratoryOpeningDate,
		"due_date":     defaultLaboratoryDueDate,
	})
	c.Equal(http.StatusCreated, status)
	laboratoryUUID := laboratoryCreationResponse["uuid"].(string)

	// Create a markdown block
	blockCreationResponse, status := CreateMarkdownBlock(cookie, laboratoryUUID)
	c.Equal(http.StatusCreated, status)
	markdownBlockUUID := blockCreationResponse["uuid"].(string)

	// Get the laboratory
	response, status := GetLaboratoryByUUID(cookie, laboratoryUUID)
	c.Equal(http.StatusOK, status)
	c.Equal(1, len(response["markdown_blocks"].([]interface{})))

	// Delete the markdown block
	_, status = DeleteMarkdownBlock(cookie, markdownBlockUUID)
	c.Equal(http.StatusNoContent, status)

	// Verify that the block was deleted
	response, status = GetLaboratoryByUUID(cookie, laboratoryUUID)
	c.Equal(http.StatusOK, status)
	c.Equal(0, len(response["markdown_blocks"].([]interface{})))
}

func TestUpdateTestBlock(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, status := CreateCourse("Update test block test - course")
	c.Equal(http.StatusCreated, status)

	// Create a laboratory
	laboratoryCreationResponse, status := CreateLaboratory(cookie, map[string]interface{}{
		"name":         "Update test block test - laboratory",
		"course_uuid":  courseUUID,
		"opening_date": defaultLaboratoryOpeningDate,
		"due_date":     defaultLaboratoryDueDate,
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

	// Update the test block
	newName := "Update test block test - block - updated"
	zipFile, err = GetSampleTestsArchive()
	c.Nil(err)

	_, status = UpdateTestBlock(&UpdateTestBlockUtilsDTO{
		blockUUID:    testBlockUUID,
		languageUUID: firstLanguageUUID,
		blockName:    newName,
		cookie:       cookie,
		testFile:     zipFile,
	})
	c.Equal(http.StatusNoContent, status)

	// Check that the test block data was updated
	laboratoryResponse, status := GetLaboratoryByUUID(cookie, laboratoryUUID)
	c.Equal(http.StatusOK, status)

	blocks := laboratoryResponse["test_blocks"].([]interface{})
	c.Equal(1, len(blocks))

	block := blocks[0].(map[string]interface{})
	c.Equal(newName, block["name"].(string))
	c.Equal(firstLanguageUUID, block["language_uuid"].(string))

	// Update the test block without sending a new test archive
	newName = "Update test block test - block - updated - 2"
	_, status = UpdateTestBlock(&UpdateTestBlockUtilsDTO{
		blockUUID:    testBlockUUID,
		languageUUID: firstLanguageUUID,
		blockName:    newName,
		cookie:       cookie,
		testFile:     nil,
	})
	c.Equal(http.StatusNoContent, status)

	// Check that the test block data was updated
	laboratoryResponse, status = GetLaboratoryByUUID(cookie, laboratoryUUID)
	c.Equal(http.StatusOK, status)

	blocks = laboratoryResponse["test_blocks"].([]interface{})
	c.Equal(1, len(blocks))

	block = blocks[0].(map[string]interface{})
	c.Equal(newName, block["name"].(string))
	c.Equal(firstLanguageUUID, block["language_uuid"].(string))
}

func TestDeleteTestBlock(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, status := CreateCourse("Delete test block test - course")
	c.Equal(http.StatusCreated, status)

	// Create a laboratory
	laboratoryCreationResponse, status := CreateLaboratory(cookie, map[string]interface{}{
		"name":         "Delete test block test - laboratory",
		"course_uuid":  courseUUID,
		"opening_date": defaultLaboratoryOpeningDate,
		"due_date":     defaultLaboratoryDueDate,
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
		blockName:      "Delete test block test - block",
		cookie:         cookie,
		testFile:       zipFile,
	})
	c.Equal(http.StatusCreated, status)
	testBlockUUID := blockCreationResponse["uuid"].(string)

	// Check that the test block was created
	laboratoryResponse, status := GetLaboratoryByUUID(cookie, laboratoryUUID)
	c.Equal(http.StatusOK, status)

	blocks := laboratoryResponse["test_blocks"].([]interface{})
	c.Equal(1, len(blocks))

	// Delete the test block
	_, status = DeleteTestBlock(cookie, testBlockUUID)
	c.Equal(http.StatusNoContent, status)

	// Check that the test block was deleted
	laboratoryResponse, status = GetLaboratoryByUUID(cookie, laboratoryUUID)
	c.Equal(http.StatusOK, status)

	blocks = laboratoryResponse["test_blocks"].([]interface{})
	c.Equal(0, len(blocks))
}

func TestSwapBlocks(t *testing.T) {
	c := require.New(t)

	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, _ := CreateCourse("Swap blocks test - course")

	// Create a laboratory
	laboratoryCreationResponse, _ := CreateLaboratory(cookie, map[string]interface{}{
		"name":         "Swap blocks test - laboratory",
		"course_uuid":  courseUUID,
		"opening_date": defaultLaboratoryOpeningDate,
		"due_date":     defaultLaboratoryDueDate,
	})
	laboratoryUUID := laboratoryCreationResponse["uuid"].(string)

	// Get the supported languages
	languagesResponse, _ := GetSupportedLanguages(cookie)
	languages := languagesResponse["languages"].([]interface{})

	firstLanguage := languages[0].(map[string]interface{})
	firstLanguageUUID := firstLanguage["uuid"].(string)

	// Create a markdown block
	blockCreationResponse, _ := CreateMarkdownBlock(cookie, laboratoryUUID)
	markdownBlockUUID := blockCreationResponse["uuid"].(string)

	// Create a test block
	zipFile, err := GetSampleTestsArchive()
	c.Nil(err)

	blockCreationResponse, _ = CreateTestBlock(&CreateTestBlockUtilsDTO{
		laboratoryUUID: laboratoryUUID,
		languageUUID:   firstLanguageUUID,
		blockName:      "Swap blocks test - block 1",
		cookie:         cookie,
		testFile:       zipFile,
	})
	testBlock1UUID := blockCreationResponse["uuid"].(string)

	// Get the laboratory
	laboratoryResponse, _ := GetLaboratoryByUUID(cookie, laboratoryUUID)
	markdownBlocks := laboratoryResponse["markdown_blocks"].([]interface{})
	testBlocks := laboratoryResponse["test_blocks"].([]interface{})

	markdownBlock := markdownBlocks[0].(map[string]interface{})
	testBlock := testBlocks[0].(map[string]interface{})

	oldMarkdownBlockIndex := markdownBlock["index"].(float64)
	oldTestBlockIndex := testBlock["index"].(float64)

	// Swap the blocks
	_, status := SwapBlocks(&SwapBlocksUtilsDTO{
		FirstBlockUUID:  markdownBlockUUID,
		SecondBlockUUID: testBlock1UUID,
		Cookie:          cookie,
	})
	c.Equal(http.StatusNoContent, status)

	// Check that the blocks were swapped
	laboratoryResponse, _ = GetLaboratoryByUUID(cookie, laboratoryUUID)

	markdownBlocks = laboratoryResponse["markdown_blocks"].([]interface{})
	testBlocks = laboratoryResponse["test_blocks"].([]interface{})

	markdownBlock = markdownBlocks[0].(map[string]interface{})
	testBlock = testBlocks[0].(map[string]interface{})

	newMarkdownBlockIndex := markdownBlock["index"].(float64)
	newTestBlockIndex := testBlock["index"].(float64)

	c.Equal(oldMarkdownBlockIndex, newTestBlockIndex)
	c.Equal(oldTestBlockIndex, newMarkdownBlockIndex)
}
