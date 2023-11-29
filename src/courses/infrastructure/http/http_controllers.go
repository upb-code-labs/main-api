package infrastructure

import (
	"net/http"

	"github.com/UPB-Code-Labs/main-api/src/courses/application"
	"github.com/UPB-Code-Labs/main-api/src/courses/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/courses/infrastructure/requests"
	"github.com/UPB-Code-Labs/main-api/src/courses/infrastructure/responses"
	"github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/gin-gonic/gin"
)

type CoursesController struct {
	UseCases *application.CoursesUseCases
}

func (controller *CoursesController) HandleCreateCourse(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")

	// Parse request body
	var request requests.CreateCourseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Validate request body
	if err := infrastructure.GetValidator().Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	// Get a random color for the course
	color, err := controller.UseCases.GetRandomColor()
	if err != nil {
		c.Error(err)
		return
	}

	// Create course
	dto := &dtos.CreateCourseDTO{
		Name:        request.Name,
		TeacherUUID: teacherUUID,
		Color:       *color,
	}

	course, err := controller.UseCases.SaveCourse(dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"uuid":  course.UUID,
		"name":  course.Name,
		"color": course.Color,
	})
}

func (controller *CoursesController) HandleGetCourse(c *gin.Context) {
	user_uuid := c.GetString("session_uuid")

	// Validate course uuid
	courseUUID := c.Param("course_uuid")
	if err := infrastructure.GetValidator().Var(courseUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid course uuid",
		})
		return
	}

	// Get course
	course, err := controller.UseCases.GetCourse(user_uuid, courseUUID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uuid":  course.UUID,
		"name":  course.Name,
		"color": course.Color,
	})
}

func (controller *CoursesController) HandleGetInvitationCode(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")

	// Validate course uuid
	courseUUID := c.Param("course_uuid")
	if err := infrastructure.GetValidator().Var(courseUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid course uuid",
		})
		return
	}

	invitationCode, err := controller.UseCases.GetInvitationCode(dtos.GetInvitationCodeDTO{
		CourseUUID:  courseUUID,
		TeacherUUID: teacherUUID,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": invitationCode,
	})
}

func (controller *CoursesController) HandleJoinCourse(c *gin.Context) {
	studentUUID := c.GetString("session_uuid")

	// Validate invitation code
	invitationCode := c.Param("invitation-code")
	if err := infrastructure.GetValidator().Var(invitationCode, "len=9"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid invitation code",
		})
		return
	}

	// Join course
	course, err := controller.UseCases.JoinCourseUsingInvitationCode(&dtos.JoinCourseUsingInvitationCodeDTO{
		StudentUUID:    studentUUID,
		InvitationCode: invitationCode,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"course": gin.H{
			"uuid":  course.UUID,
			"name":  course.Name,
			"color": course.Color,
		},
	})
}

func (controller *CoursesController) HandleGetEnrolledCourses(c *gin.Context) {
	userUUID := c.GetString("session_uuid")

	// Get enrolled courses
	enrolledCourses, err := controller.UseCases.GetEnrolledCourses(userUUID)
	if err != nil {
		c.Error(err)
		return
	}

	// Parse enrolled courses to response
	c.JSON(http.StatusOK, responses.GetResponseFromDTO(enrolledCourses))
}

func (controller *CoursesController) HandleChangeCourseVisibility(c *gin.Context) {
	userUUID := c.GetString("session_uuid")

	// Validate course uuid
	courseUUID := c.Param("course_uuid")
	if err := infrastructure.GetValidator().Var(courseUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid course uuid",
		})
		return
	}

	// Change course visibility
	isHiddenAfterUpdate, err := controller.UseCases.ToggleCourseVisibility(courseUUID, userUUID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"visible": !isHiddenAfterUpdate,
	})
}

func (controller *CoursesController) HandleChangeCourseName(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")

	// Parse request body
	var request requests.CreateCourseRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Validate request body
	if err := infrastructure.GetValidator().Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	// Validate course uuid
	courseUUID := c.Param("course_uuid")
	if err := infrastructure.GetValidator().Var(courseUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid course uuid",
		})
		return
	}

	// Change course name
	err := controller.UseCases.UpdateCourseName(dtos.RenameCourseDTO{
		TeacherUUID: teacherUUID,
		CourseUUID:  courseUUID,
		NewName:     request.Name,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (controller *CoursesController) HandleAddStudentToCourse(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")

	// Validate course uuid
	courseUUID := c.Param("course_uuid")
	if err := infrastructure.GetValidator().Var(courseUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid course uuid",
		})
		return
	}

	// Parse request body
	var request requests.EnrollStudentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Validate request body
	if err := infrastructure.GetValidator().Struct(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Validation error",
			"errors":  err.Error(),
		})
		return
	}

	// Add student to course
	dto := &dtos.AddStudentToCourseDTO{
		TeacherUUID: teacherUUID,
		StudentUUID: request.StudentUUID,
		CourseUUID:  courseUUID,
	}

	err := controller.UseCases.AddStudentToCourse(dto)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (controller *CoursesController) HandleGetStudentsEnrolledInCourse(c *gin.Context) {
	teacherUUID := c.GetString("session_uuid")

	// Validate course uuid
	courseUUID := c.Param("course_uuid")
	if err := infrastructure.GetValidator().Var(courseUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not valid course uuid",
		})
		return
	}

	// Get enrolled students
	enrolledStudents, err := controller.UseCases.GetEnrolledStudents(teacherUUID, courseUUID)
	if err != nil {
		c.Error(err)
		return
	}

	enrolledStudentsResponse := make([]gin.H, len(enrolledStudents))
	for i, student := range enrolledStudents {
		enrolledStudentsResponse[i] = gin.H{
			"uuid":             student.UUID,
			"full_name":        student.FullName,
			"institutional_id": student.InstitutionalId,
			"is_active":        student.IsActive,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"students": enrolledStudentsResponse,
	})
}

func (controller *CoursesController) HandleGetCourseLaboratories(c *gin.Context) {
	userUUID := c.GetString("session_uuid")
	userRole := c.GetString("session_role")

	// Validate course uuid
	courseUUID := c.Param("course_uuid")
	if err := infrastructure.GetValidator().Var(courseUUID, "uuid4"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Not valid course uuid",
		})
		return
	}

	// Get course laboratories
	laboratories, err := controller.UseCases.GetCourseLaboratories(dtos.GetCourseLaboratoriesDTO{
		CourseUUID: courseUUID,
		UserUUID:   userUUID,
		UserRole:   userRole,
	})
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"laboratories": laboratories,
	})
}
