package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	accounts_application "github.com/UPB-Code-Labs/main-api/src/accounts/application"
	accounts_infrastructure "github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure"
	"github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure/requests"
	config_infrastructure "github.com/UPB-Code-Labs/main-api/src/config/infrastructure"
	session_application "github.com/UPB-Code-Labs/main-api/src/session/application"
	session_infrastructure "github.com/UPB-Code-Labs/main-api/src/session/infrastructure"
	shared_infrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

// --- Globals ---
var (
	router              *gin.Engine
	accountsControllers *accounts_infrastructure.AccountsController
	sessionControllers  *session_infrastructure.SessionControllers

	registeredStudentEmail string
	registeredStudentPass  string

	registeredAdminEmail = "development.admin@gmail.com"
	registeredAdminPass  = "changeme123*/"

	registeredTeacherEmail string
	registeredTeacherPass  string
)

// --- Setup ---
func TestMain(m *testing.M) {
	// Setup
	setupDatabase()
	setupRouter()
	setupControllers()
	registerRoutes()
	registerBaseAccounts()
	defer shared_infrastructure.ClosePostgresConnection()

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func setupDatabase() {
	shared_infrastructure.GetPostgresConnection()
	config_infrastructure.RunMigrations()
}

func setupRouter() {
	router = gin.Default()
	router.Use(shared_infrastructure.ErrorHandlerMiddleware())
}

func setupControllers() {
	setupAccountsControllers()
	setupSessionControllers()
}

func setupAccountsControllers() {
	useCases := accounts_application.AccountsUseCases{
		AccountsRepository: accounts_infrastructure.GetAccountsPgRepository(),
		PasswordsHasher:    accounts_infrastructure.GetArgon2PasswordsHasher(),
	}

	controllers := &accounts_infrastructure.AccountsController{
		UseCases: &useCases,
	}

	accountsControllers = controllers
}

func setupSessionControllers() {
	useCases := session_application.SessionUseCases{
		AccountsRepository: accounts_infrastructure.GetAccountsPgRepository(),
		PasswordHasher:     accounts_infrastructure.GetArgon2PasswordsHasher(),
		TokenHandler:       shared_infrastructure.GetJwtTokenHandler(),
	}

	controllers := &session_infrastructure.SessionControllers{
		UseCases: &useCases,
	}

	sessionControllers = controllers
}

func registerBaseAccounts() {
	registerBaseStudent()
	registerBaseTeacher()
}

func registerBaseStudent() {
	studentEmail := "greta.mann.2020@upb.edu.co"
	studentPassword := "greta/password/2023"

	code := RegisterStudent(requests.RegisterUserRequest{
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

func registerBaseTeacher() {
	teacherEmail := "judy.arroyo.2020@upb.edu.co"
	teacherPassword := "judy/password/2023"

	code := RegisterTeacherWithoutAuth(requests.RegisterTeacherRequest{
		FullName: "Judy Arroyo",
		Email:    teacherEmail,
		Password: teacherPassword,
	})
	if code != http.StatusCreated {
		panic("Error registering base teacher")
	}

	registeredTeacherEmail = teacherEmail
	registeredTeacherPass = teacherPassword
}

func registerRoutes() {
	// Session
	router.POST("/session/login", sessionControllers.HandleLogin)
	router.DELETE(
		"/session/logout",
		shared_infrastructure.WithAuthenticationMiddleware(),
		sessionControllers.HandleLogout,
	)
	router.GET(
		"/session/whoami",
		shared_infrastructure.WithAuthenticationMiddleware(),
		sessionControllers.HandleWhoAmI,
	)

	// Accounts
	router.POST("/accounts/students", accountsControllers.HandleRegisterStudent)

	router.POST("/accounts/admins/no_auth", accountsControllers.HandleRegisterAdmin)
	router.POST(
		"/accounts/admins",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware("admin"),
		accountsControllers.HandleRegisterAdmin,
	)

	router.POST("/accounts/teachers/no_auth", accountsControllers.HandleRegisterTeacher)
	router.POST(
		"/accounts/teachers",
		shared_infrastructure.WithAuthenticationMiddleware(),
		shared_infrastructure.WithAuthorizationMiddleware("admin"),
		accountsControllers.HandleRegisterTeacher,
	)
}

// --- Helpers ---
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

func ParseJsonResponse(buffer *bytes.Buffer) map[string]interface{} {
	var response map[string]interface{}
	json.Unmarshal(buffer.Bytes(), &response)
	return response
}
