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
			"message": "Invalid request body",
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
