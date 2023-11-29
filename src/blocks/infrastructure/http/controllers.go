package http

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/blocks/application"
	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/blocks/infrastructure/requests"
	shared_infrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

type BlocksController struct {
	UseCases *application.BlocksUseCases
}

func (controller *BlocksController) HandleUpdateMarkdownBlockContent(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")
	blockUUID := c.Param("block_uuid")

	// Validate the block UUID
	if err := shared_infrastructure.GetValidator().Var(blockUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Block UUID is not valid",
		})
		return
	}

	// Parse request body
	var request requests.UpdateMarkdownBlockContentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Request body is not valid",
		})
		return
	}

	// Validate request body
	if err := shared_infrastructure.GetValidator().Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	dto := dtos.UpdateMarkdownBlockContentDTO{
		BlockUUID:   blockUUID,
		TeacherUUID: teacherUUID,
		Content:     request.Content,
	}

	err := controller.UseCases.UpdateMarkdownBlockContent(dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
