package http

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/blocks/application"
	"github.com/gin-gonic/gin"
)

type BlocksController struct {
	UseCases *application.BlocksUseCases
}

func (controller *BlocksController) HandleUpdateMarkdownBlockContent(c *gin.Context) {
	c.Status(http.StatusNotImplemented)
}
