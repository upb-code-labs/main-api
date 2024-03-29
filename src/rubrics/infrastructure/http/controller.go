package http

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/rubrics/application"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/infrastructure/requests"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
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
	if err := sharedInfrastructure.GetValidator().Struct(request); err != nil {
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
	// Validate rubric UUID
	rubric_uuid := c.Param("rubricUUID")
	if err := sharedInfrastructure.GetValidator().Var(rubric_uuid, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid rubric uuid",
		})
		return
	}

	// Get the rubric
	rubric, err := controller.UseCases.GetRubricByUUID(rubric_uuid)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rubric": rubric,
	})
}

func (controller *RubricsController) HandleDeleteRubric(c *gin.Context) {
	teacher_uuid := c.GetString("session_uuid")

	// Validate rubric UUID
	rubric_uuid := c.Param("rubricUUID")
	if err := sharedInfrastructure.GetValidator().Var(rubric_uuid, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid rubric uuid",
		})
		return
	}

	// Create DTO
	dto := dtos.DeleteRubricDTO{
		TeacherUUID: teacher_uuid,
		RubricUUID:  rubric_uuid,
	}

	// Delete the rubric
	err := controller.UseCases.DeleteRubric(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (controller *RubricsController) HandleUpdateRubricName(c *gin.Context) {
	teacher_uuid := c.GetString("session_uuid")

	// Validate rubric UUID
	rubric_uuid := c.Param("rubricUUID")
	if err := sharedInfrastructure.GetValidator().Var(rubric_uuid, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid rubric uuid",
		})
		return
	}

	// Parse request body
	var request requests.UpdateRubricNameRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
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

	// Create DTO
	dto := dtos.UpdateRubricNameDTO{
		TeacherUUID: teacher_uuid,
		RubricUUID:  rubric_uuid,
		Name:        request.Name,
	}

	// Update the rubric
	err := controller.UseCases.UpdateRubricName(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (controller *RubricsController) HandleAddObjectiveToRubric(c *gin.Context) {
	teacher_uuid := c.GetString("session_uuid")

	// Validate rubric UUID
	rubric_uuid := c.Param("rubricUUID")
	if err := sharedInfrastructure.GetValidator().Var(rubric_uuid, "uuid4"); err != nil {
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
	if err := sharedInfrastructure.GetValidator().Struct(request); err != nil {
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
	if err := sharedInfrastructure.GetValidator().Var(objective_uuid, "uuid4"); err != nil {
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
	if err := sharedInfrastructure.GetValidator().Struct(request); err != nil {
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
	if err := sharedInfrastructure.GetValidator().Var(objective_uuid, "uuid4"); err != nil {
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
	if err := sharedInfrastructure.GetValidator().Struct(request); err != nil {
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

func (controller *RubricsController) HandleDeleteObjective(c *gin.Context) {
	teacher_uuid := c.GetString("session_uuid")

	// Validate objective UUID
	objective_uuid := c.Param("objectiveUUID")
	if err := sharedInfrastructure.GetValidator().Var(objective_uuid, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid objective uuid",
		})
		return
	}

	// Create DTO
	dto := dtos.DeleteObjectiveDTO{
		TeacherUUID:   teacher_uuid,
		ObjectiveUUID: objective_uuid,
	}

	// Delete the objective
	err := controller.UseCases.DeleteObjective(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (controller *RubricsController) HandleUpdateCriteria(c *gin.Context) {
	teacher_uuid := c.GetString("session_uuid")

	// Validate criteria UUID
	criteria_uuid := c.Param("criteriaUUID")
	if err := sharedInfrastructure.GetValidator().Var(criteria_uuid, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid criteria uuid",
		})
		return
	}

	// Parse request body
	var request requests.UpdateCriteriaRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
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

	// Create DTO
	dto := dtos.UpdateCriteriaDTO{
		TeacherUUID:         teacher_uuid,
		CriteriaUUID:        criteria_uuid,
		CriteriaDescription: request.Description,
		CriteriaWeight:      request.Weight,
	}

	// Update the criteria
	err := controller.UseCases.UpdateCriteria(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (controller *RubricsController) HandleDeleteCriteria(c *gin.Context) {
	teacher_uuid := c.GetString("session_uuid")

	// Validate criteria UUID
	criteria_uuid := c.Param("criteriaUUID")
	if err := sharedInfrastructure.GetValidator().Var(criteria_uuid, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid criteria uuid",
		})
		return
	}

	// Create DTO
	dto := dtos.DeleteCriteriaDTO{
		TeacherUUID:  teacher_uuid,
		CriteriaUUID: criteria_uuid,
	}

	// Delete the criteria
	err := controller.UseCases.DeleteCriteria(&dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
