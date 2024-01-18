package http

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/grades/application"
	"github.com/UPB-Code-Labs/main-api/src/grades/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/grades/infrastructure/requests"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

// GradesController controller to handle the requests to the `/grades` endpoints
type GradesController struct {
	UseCases *application.GradesUseCases
}

// HandleGetSummarizedGradesInLaboratory controller to get the summarized grades of the students in a laboratory
func (controller *GradesController) HandleGetSummarizedGradesInLaboratory(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")
	laboratoryUUID := c.Param("laboratoryUUID")

	// Validate laboratory UUID
	if err := sharedInfrastructure.GetValidator().Var(laboratoryUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Laboratory UUID is not valid",
		})
		return
	}

	// Get summarized grades
	summarizedGrades, err := controller.UseCases.GetSummarizedGradesInLaboratory(&dtos.GetSummarizedGradesInLaboratoryDTO{
		TeacherUUID:    teacherUUID,
		LaboratoryUUID: laboratoryUUID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"grades": summarizedGrades,
	})
}

// HandleSetCriteriaGrade controller to select a criteria from a rubric to be added to a student's grade
func (controller *GradesController) HandleSetCriteriaGrade(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")
	laboratoryUUID := c.Param("laboratoryUUID")
	studentUUID := c.Param("studentUUID")

	// Validate UUIDs
	uuids := []string{teacherUUID, laboratoryUUID, studentUUID}
	for _, uuid := range uuids {
		if err := sharedInfrastructure.GetValidator().Var(uuid, "uuid4"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Please, make sure you are sending valid UUIDs",
			})
			return
		}
	}

	// Parse the request body
	var request requests.SetCriteriaToGradeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Request body is not valid",
		})
		return
	}

	// Validate the request body
	if err := sharedInfrastructure.GetValidator().Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	// Create DTO, note that the rubric field will be populated in the use case
	dto := &dtos.SetCriteriaToGradeDTO{
		TeacherUUID:    teacherUUID,
		LaboratoryUUID: laboratoryUUID,
		StudentUUID:    studentUUID,
		ObjectiveUUID:  request.ObjectiveUUID,
		CriteriaUUID:   request.CriteriaUUID,
	}

	// Set criteria to grade
	err := controller.UseCases.SetCriteriaToGrade(dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (controller *GradesController) HandleGetStudentGradeInLaboratoryWithRubric(c *gin.Context) {
	// Get user uuid from the session
	userUUID := c.GetString("session_uuid")

	// Get UUIDs from the request params
	studentUUID := c.Param("studentUUID")
	laboratoryUUID := c.Param("laboratoryUUID")
	rubricUUID := c.Param("rubricUUID")

	// Validate UUIDs
	uuids := []string{userUUID, studentUUID, laboratoryUUID, rubricUUID}
	for _, uuid := range uuids {
		if err := sharedInfrastructure.GetValidator().Var(uuid, "uuid4"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Please, make sure you are sending valid UUIDs",
			})
			return
		}
	}

	// Get student grade in laboratory with rubric
	grade, err := controller.UseCases.GetStudentGradeInLaboratoryWithRubric(
		&dtos.GetStudentGradeInLaboratoryWithRubricDTO{
			UserUUID:       userUUID,
			StudentUUID:    studentUUID,
			LaboratoryUUID: laboratoryUUID,
			RubricUUID:     rubricUUID,
		},
	)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, grade)
}

// HandleSetCommentToGrade controller to set a comment to a student's grade
func (controller *GradesController) HandleSetCommentToGrade(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")
	laboratoryUUID := c.Param("laboratoryUUID")
	studentUUID := c.Param("studentUUID")

	// Validate UUIDs
	uuids := []string{teacherUUID, laboratoryUUID, studentUUID}
	for _, uuid := range uuids {
		if err := sharedInfrastructure.GetValidator().Var(uuid, "uuid4"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Please, make sure you are sending valid UUIDs",
			})
			return
		}
	}

	// Parse the request body
	var request requests.SetCommentToGradeRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Request body is not valid",
		})
		return
	}

	// Validate the request body
	if err := sharedInfrastructure.GetValidator().Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	// Create DTO. Note that the rubric field will be populated in the use case
	dto := &dtos.SetCommentToGradeDTO{
		TeacherUUID:    teacherUUID,
		LaboratoryUUID: laboratoryUUID,
		StudentUUID:    studentUUID,
		Comment:        request.Comment,
	}

	err := controller.UseCases.SetCommentToGrade(dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}
