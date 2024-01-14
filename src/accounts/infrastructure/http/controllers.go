package infrastructure

import (
	"errors"
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/accounts/application"
	"github.com/UPB-Code-Labs/main-api/src/accounts/domain/dtos"
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
			"message": "Request body is not valid",
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
			"message": "Request body is not valid",
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
			"message": "Request body is not valid",
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

func (controller *AccountsController) HandleUpdatePassword(c *gin.Context) {
	userUUID := c.GetString("session_uuid")

	// Parse request body
	var request requests.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{
			"message": "Request body is not valid",
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

	// Update password
	dto := dtos.UpdatePasswordDTO{
		UserUUID:    userUUID,
		OldPassword: request.OldPassword,
		NewPassword: request.NewPassword,
	}

	err := controller.UseCases.UpdatePassword(dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// HandleUpdateProfile controller to update the profile of an account
func (controller *AccountsController) HandleUpdateProfile(c *gin.Context) {
	userUUID := c.GetString("session_uuid")
	userRole := c.GetString("session_role")

	// Parse and validate the request according to the user role
	var request interface{}
	switch userRole {
	case "admin":
		request = &requests.UpdateAdminProfileRequest{}
	case "teacher":
		request = &requests.UpdateTeacherProfileRequest{}
	default:
		request = &requests.UpdateStudentProfileRequest{}
	}

	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(400, gin.H{
			"message": "Request body is not valid",
		})
		return
	}

	if err := infrastructure.GetValidator().Struct(request); err != nil {
		c.JSON(400, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	// Create the DTO to update the profile
	var fullName, email, password string
	var institutionalId *string

	switch v := request.(type) {
	case *requests.UpdateAdminProfileRequest:
		fullName = v.FullName
		email = v.Email
		password = v.Password
	case *requests.UpdateTeacherProfileRequest:
		fullName = v.FullName
		email = v.Email
		password = v.Password
	case *requests.UpdateStudentProfileRequest:
		fullName = v.FullName
		email = v.Email
		institutionalId = v.InstitutionalId
		password = v.Password
	default:
		c.Error(errors.New("request type not supported"))
	}

	updateAccountDTO := dtos.UpdateAccountDTO{
		UserUUID:        userUUID,
		FullName:        fullName,
		Email:           email,
		InstitutionalId: institutionalId,
		Password:        password,
	}

	// Update profile
	err := controller.UseCases.UpdateProfile(updateAccountDTO)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
