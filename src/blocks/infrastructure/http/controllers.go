package http

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/blocks/application"
	"github.com/UPB-Code-Labs/main-api/src/blocks/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/blocks/infrastructure/requests"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

type BlocksController struct {
	UseCases *application.BlocksUseCases
}

func (controller *BlocksController) HandleUpdateMarkdownBlockContent(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")
	blockUUID := c.Param("block_uuid")

	// Validate the block UUID
	if err := sharedInfrastructure.GetValidator().Var(blockUUID, "uuid4"); err != nil {
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
	if err := sharedInfrastructure.GetValidator().Struct(request); err != nil {
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

func (controller *BlocksController) HandleUpdateTestBlock(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")
	blockUUID := c.Param("block_uuid")

	// Validate the request struct
	languageUUID := c.PostForm("language_uuid")
	blockName := c.PostForm("block_name")

	req := requests.UpdateTestBlockRequest{
		LanguageUUID: languageUUID,
		Name:         blockName,
	}

	if err := sharedInfrastructure.GetValidator().Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	// Create the DTO
	dto := dtos.UpdateTestBlockDTO{
		TeacherUUID:  teacherUUID,
		BlockUUID:    blockUUID,
		LanguageUUID: languageUUID,
		Name:         blockName,
	}

	// Validate the test archive (if any)
	multipartHeader, err := c.FormFile("test_archive")
	if err != nil {
		if err != http.ErrMissingFile {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Please, make sure to send the test archive",
			})
			return
		}
	}

	if multipartHeader != nil {
		err = sharedInfrastructure.ValidateMultipartFileHeader(multipartHeader)
		if err != nil {
			c.Error(err)
			return
		}

		// Add the test archive to the DTO
		multipartFile, err := multipartHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "There was an error while reading the test archive",
			})
			return
		}

		dto.NewTestArchive = &multipartFile
	}

	// Update the test block
	err = controller.UseCases.UpdateTestBlock(dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (controller *BlocksController) HandleDeleteMarkdownBlock(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")
	blockUUID := c.Param("block_uuid")

	// Validate the block UUID
	if err := sharedInfrastructure.GetValidator().Var(blockUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Block UUID is not valid",
		})
		return
	}

	// Delete the block
	err := controller.UseCases.DeleteMarkdownBlock(dtos.DeleteBlockDTO{
		TeacherUUID: teacherUUID,
		BlockUUID:   blockUUID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (controller *BlocksController) HandleDeleteTestBlock(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")
	blockUUID := c.Param("block_uuid")

	// Validate the block UUID
	if err := sharedInfrastructure.GetValidator().Var(blockUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Block UUID is not valid",
		})
		return
	}

	// Delete the block
	err := controller.UseCases.DeleteTestBlock(dtos.DeleteBlockDTO{
		TeacherUUID: teacherUUID,
		BlockUUID:   blockUUID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (controller *BlocksController) HandleSwapBlocks(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")

	// Parse request body
	var request requests.SwapBlocksRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Request body is not valid",
		})
		return
	}

	// Validate request body
	if err := sharedInfrastructure.GetValidator().Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	dto := dtos.SwapBlocksDTO{
		TeacherUUID:     teacherUUID,
		FirstBlockUUID:  request.FirstBlockUUID,
		SecondBlockUUID: request.SecondBlockUUID,
	}

	err := controller.UseCases.SwapBlocks(dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

// HandleGetTestBlockTestsArchive controller to handle the request of downloading the `.zip` archive
// containing the tests of a test block
func (controller *BlocksController) HandleGetTestBlockTestsArchive(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")
	blockUUID := c.Param("block_uuid")

	// Validate the block UUID
	if err := sharedInfrastructure.GetValidator().Var(blockUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Block UUID is not valid",
		})
		return
	}

	// Get the test block tests archive
	testsArchive, err := controller.UseCases.GetTestBlockTestsArchive(&dtos.GetBlockTestsArchiveDTO{
		TeacherUUID: teacherUUID,
		BlockUUID:   blockUUID,
	})

	if err != nil {
		c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/zip", testsArchive)
}
