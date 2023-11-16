package http

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/rubrics/application"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/infrastructure/requests"
	shared_infrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

type RubricsController struct {
	UseCases *application.RubricsUseCases
}

func (controller *RubricsController) HandleCreateRubric(c *gin.Context) {
	teacher_uuid := c.GetString("session_uuid")

	// Parse request body
	var request requests.CreateRubricRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
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

	// Create DTO
	dto := dtos.CreateRubricDTO{
		TeacherUUID: teacher_uuid,
		Name:        request.Name,
	}

	// Create the course
	rubric, err := controller.UseCases.CreateRubric(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Rubric created",
		"uuid":    rubric.UUID,
		"name":    rubric.Name,
	})
}

func (controller *RubricsController) HandleGetRubricsCreatedByTeacher(c *gin.Context) {
	teacher_uuid := c.GetString("session_uuid")

	rubrics, err := controller.UseCases.GetRubricsCreatedByTeacher(teacher_uuid)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rubrics": rubrics,
		"message": "Rubrics were retrieved",
	})
}

func (controller *RubricsController) HandleGetRubricByUUID(c *gin.Context) {
	teacher_uuid := c.GetString("session_uuid")

	// Validate rubric UUID
	rubric_uuid := c.Param("rubricUUID")
	if err := shared_infrastructure.GetValidator().Var(rubric_uuid, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid rubric uuid",
		})
		return
	}

	// Create DTO
	dto := dtos.GetRubricDto{
		TeacherUUID: teacher_uuid,
		RubricUUID:  rubric_uuid,
	}

	// Get the rubric
	rubric, err := controller.UseCases.GetRubricByUUID(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rubric was retrieved",
		"rubric":  rubric,
	})
}

func (controller *RubricsController) HandleAddObjectiveToRubric(c *gin.Context) {
	teacher_uuid := c.GetString("session_uuid")

	// Validate rubric UUID
	rubric_uuid := c.Param("rubricUUID")
	if err := shared_infrastructure.GetValidator().Var(rubric_uuid, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid rubric uuid",
		})
		return
	}

	// Parse request body
	var request requests.AddObjectiveToRubricRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
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

	// Create DTO
	dto := dtos.AddObjectiveToRubricDTO{
		TeacherUUID:          teacher_uuid,
		RubricUUID:           rubric_uuid,
		ObjectiveDescription: request.Description,
	}

	// Add the objective
	objective_uuid, err := controller.UseCases.AddObjectiveToRubric(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Objective added to rubric",
		"uuid":    objective_uuid,
	})
}

func (controller *RubricsController) HandleAddCriteriaToObjective(c *gin.Context) {
	teacher_uuid := c.GetString("session_uuid")

	// Validate objective UUID
	objective_uuid := c.Param("objectiveUUID")
	if err := shared_infrastructure.GetValidator().Var(objective_uuid, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid objective uuid",
		})
		return
	}

	// Parse request body
	var request requests.AddCriteriaToObjectiveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
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

	// Create DTO
	dto := dtos.AddCriteriaToObjectiveDTO{
		TeacherUUID:         teacher_uuid,
		ObjectiveUUID:       objective_uuid,
		CriteriaDescription: request.Description,
		CriteriaWeight:      request.Weight,
	}

	// Add the criteria
	criteria_uuid, err := controller.UseCases.AddCriteriaToObjective(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Criteria added to objective",
		"uuid":    criteria_uuid,
	})
}

func (controller *RubricsController) HandleUpdateObjective(c *gin.Context) {
	teacher_uuid := c.GetString("session_uuid")

	// Validate objective UUID
	objective_uuid := c.Param("objectiveUUID")
	if err := shared_infrastructure.GetValidator().Var(objective_uuid, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid objective uuid",
		})
		return
	}

	// Parse request body
	var request requests.UpdateObjectiveRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
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

	// Create DTO
	dto := dtos.UpdateObjectiveDTO{
		TeacherUUID:        teacher_uuid,
		ObjectiveUUID:      objective_uuid,
		UpdatedDescription: request.Description,
	}

	// Update the objective
	err := controller.UseCases.UpdateObjective(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
