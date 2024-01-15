package http

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/grades/application"
	"github.com/UPB-Code-Labs/main-api/src/grades/domain/dtos"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

// GradesController controller to handle the requests to the `/grades` endpoints
type GradesController struct {
	UseCases *application.GradesUseCases
}

// GetSummarizedGradesInLaboratory controller to get the summarized grades of the students in a laboratory
func (controller *GradesController) GetSummarizedGradesInLaboratory(c *gin.Context) {
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
