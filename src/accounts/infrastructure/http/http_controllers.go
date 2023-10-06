package infrastructure

import (
	"github.com/UPB-Code-Labs/main-api/src/accounts/application"
	"github.com/UPB-Code-Labs/main-api/src/accounts/infrastructure/requests"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

type AccountsController struct {
	UseCases *application.AccountsUseCases
}

func (controller *AccountsController) HandleRegisterStudent(c *gin.Context) {
	// Parse request body
	var request requests.RegisterUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Validate request body
	if err := infrastructure.GetValidator().Struct(request); err != nil {
		c.JSON(400, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	// Register student
	dto := request.ToDTO()
	err := controller.UseCases.RegisterStudent(*dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(201)
}

func (controller *AccountsController) HandleRegisterAdmin(c *gin.Context) {
	adminUUID := c.GetString("session_uuid")

	// Parse request body
	var request requests.RegisterAdminRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Validate request body
	if err := infrastructure.GetValidator().Struct(request); err != nil {
		c.JSON(400, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	// Register admin
	dto := request.ToDTO()
	dto.CreatedBy = adminUUID
	err := controller.UseCases.RegisterAdmin(*dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(201)
}

func (controller *AccountsController) HandleRegisterTeacher(c *gin.Context) {
	adminUUID := c.GetString("session_uuid")

	// Parse request body
	var request requests.RegisterTeacherRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Validate request body
	if err := infrastructure.GetValidator().Struct(request); err != nil {
		c.JSON(400, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	// Register teacher
	dto := request.ToDTO()
	dto.CreatedBy = adminUUID
	err := controller.UseCases.RegisterTeacher(*dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(201)
}

func (controller *AccountsController) HandleGetAdmins(c *gin.Context) {
	admins, err := controller.UseCases.GetAdmins()
	if err != nil {
		c.Error(err)
		return
	}

	publicAdmins := make([]gin.H, len(admins))
	for i, admin := range admins {
		publicAdmins[i] = gin.H{
			"uuid":       admin.UUID,
			"full_name":  admin.FullName,
			"created_at": admin.CreatedAt,
			"created_by": admin.CreatedBy,
		}
	}

	c.JSON(200, gin.H{
		"admins": publicAdmins,
	})
}

func (controller *AccountsController) HandleSearchStudents(c *gin.Context) {
	// Get query params
	fullName := c.Query("fullName")
	if fullName == "" {
		c.JSON(400, gin.H{
			"message": "Missing fullName query param",
		})
		return
	}

	// Search students
	students, err := controller.UseCases.SearchStudentsByFullName(fullName)
	if err != nil {
		c.Error(err)
		return
	}

	publicStudents := make([]gin.H, len(students))
	for i, student := range students {
		publicStudents[i] = gin.H{
			"uuid":             student.UUID,
			"full_name":        student.FullName,
			"institutional_id": student.InstitutionalId,
		}
	}

	c.JSON(200, gin.H{
		"students": publicStudents,
	})
}
