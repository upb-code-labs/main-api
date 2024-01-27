package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure/requests"
	configInfrastructure "github.com/UPB-Code-Labs/main-api/src/config/infrastructure"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	submissionsImplementations "github.com/UPB-Code-Labs/main-api/src/submissions/infrastructure/implementations"
	"github.com/gin-gonic/gin"
)

// --- Globals ---
var (
	router *gin.Engine

	registeredStudentEmail string
	registeredStudentPass  string

	registeredAdminEmail = "development.admin@gmail.com"
	registeredAdminPass  = "changeme123*/"

	registeredTeacherEmail string
	registeredTeacherPass  string

	secondRegisteredTeacherEmail string
	secondRegisteredTeacherPass  string

	defaultLaboratoryOpeningDate = "2023-12-01T12:00:00-05:00"
	defaultLaboratoryDueDate     = "3023-12-01T12:00:00-05:00"

	defaultLaboratoryOpeningDateUTC = "2023-12-01T17:00:00Z"
	defaultLaboratoryDueDateUTC     = "3023-12-01T17:00:00Z"
)

type GenericTestCase struct {
	Payload            map[string]interface{}
	ExpectedStatusCode int
}

// --- Setup ---
func TestMain(m *testing.M) {
	// Setup database
	setupDatabase()
	defer sharedInfrastructure.ClosePostgresConnection()

	// Setup RabbitMQ
	setupRabbitMQ()
	defer sharedInfrastructure.CloseRabbitMQConnection()

	// Setup SSE
	setupSSE()

	// Setup http router
	setupRouter()
	registerBaseAccounts()

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func setupDatabase() {
	sharedInfrastructure.GetPostgresConnection()
	configInfrastructure.RunMigrations()
}

func setupRabbitMQ() {
	// Connect to RabbitMQ
	sharedInfrastructure.ConnectToRabbitMQ()

	// Start listening for messages in the submissions real time updates queue
	submissionsRealTimeUpdatesQueueMgr := submissionsImplementations.GetSubmissionsRealTimeUpdatesQueueMgrInstance()
	go submissionsRealTimeUpdatesQueueMgr.ListenForUpdates()
}

func setupSSE() {
	// Start listening for SSE connections
	realTimeSubmissionsUpdatesSender := submissionsImplementations.GetSubmissionsRealTimeUpdatesSenderInstance()
	go realTimeSubmissionsUpdatesSender.Listen()
}

func setupRouter() {
	router = configInfrastructure.InstanceHttpServer()
}

func registerBaseAccounts() {
	registerBaseStudent()
	registerBaseTeachers()
}

func registerBaseStudent() {
	studentEmail := "greta.mann.2020@upb.edu.co"
	studentPassword := "greta/password/2023"

	code := RegisterStudentAccount(requests.RegisterUserRequest{
		FullName:        "Greta Mann",
		Email:           studentEmail,
		InstitutionalId: "000123456",
		Password:        studentPassword,
	})
	if code != http.StatusCreated {
		panic("Error registering base student")
	}

	registeredStudentEmail = studentEmail
	registeredStudentPass = studentPassword
}

func registerBaseTeachers() {
	// Register the first teacher
	teacherEmail := "judy.arroyo.2020@upb.edu.co"
	teacherPassword := "judy/password/2023"

	code := RegisterTeacherAccount(requests.RegisterTeacherRequest{
		FullName: "Judy Arroyo",
		Email:    teacherEmail,
		Password: teacherPassword,
	})
	if code != http.StatusCreated {
		panic("Error registering base teacher")
	}

	registeredTeacherEmail = teacherEmail
	registeredTeacherPass = teacherPassword

	// Register the second teacher
	secondTeacherEmail := "trofim.vijay.2020@upb.edu.co"
	secondTeacherPassword := "trofim/password/2023"

	code = RegisterTeacherAccount(requests.RegisterTeacherRequest{
		FullName: "Trofim Vijay",
		Email:    secondTeacherEmail,
		Password: secondTeacherPassword,
	})
	if code != http.StatusCreated {
		panic("Error registering base teacher")
	}

	secondRegisteredTeacherEmail = secondTeacherEmail
	secondRegisteredTeacherPass = secondTeacherPassword
}

// --- Helpers ---
func GetSampleTestsArchive() (*os.File, error) {
	TEST_FILE_PATH := "../data/java_tests_sample.zip"

	zipFile, err := os.Open(TEST_FILE_PATH)
	if err != nil {
		return nil, err
	}

	return zipFile, nil
}

func GetSampleSubmissionArchive() (*os.File, error) {
	SUBMISSION_FILE_PATH := "../data/java_submission_sample.zip"

	zipFile, err := os.Open(SUBMISSION_FILE_PATH)
	if err != nil {
		return nil, err
	}

	return zipFile, nil
}

func PrepareRequest(method, endpoint string, payload interface{}) (*httptest.ResponseRecorder, *http.Request) {
	var req *http.Request

	if payload != nil {
		payloadBytes, _ := json.Marshal(payload)
		req, _ = http.NewRequest(method, endpoint, bytes.NewReader(payloadBytes))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, _ = http.NewRequest(method, endpoint, nil)
	}

	w := httptest.NewRecorder()
	return w, req
}

func PrepareMultipartRequest(method, endpoint string, body *bytes.Buffer) (*httptest.ResponseRecorder, *http.Request) {
	var req *http.Request

	req, _ = http.NewRequest(method, endpoint, body)
	req.Header.Set("Content-Type", "multipart/form-data")

	w := httptest.NewRecorder()
	return w, req
}

func ParseJsonResponse(buffer *bytes.Buffer) map[string]interface{} {
	var response map[string]interface{}
	json.Unmarshal(buffer.Bytes(), &response)
	return response
}
