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
	config_infrastructure "github.com/UPB-Code-Labs/main-api/src/config/infrastructure"
	shared_infrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

// --- Globals ---
var (
	router              *gin.Engine
	accountsControllers *accounts_infrastructure.AccountsController
)

type TestCase struct {
	Payload            map[string]interface{}
	ExpectedStatusCode int
}

// --- Setup ---
func TestMain(m *testing.M) {
	setupDatabase()
	defer shared_infrastructure.ClosePostgresConnection()
	setupRouter()
	setupControllers()

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

// --- Helpers ---
func PerformRequest(method, endpoint string, payload interface{}) (*httptest.ResponseRecorder, *http.Request) {
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
