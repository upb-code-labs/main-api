package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetCriteriaToStudentSubmission(t *testing.T) {
	c := require.New(t)

	// ## Test preparation
	// Login as a teacher
	w, r := PrepareRequest("POST", "/api/v1/session/login", map[string]interface{}{
		"email":    registeredTeacherEmail,
		"password": registeredTeacherPass,
	})
	router.ServeHTTP(w, r)
	cookie := w.Result().Cookies()[0]

	// Create a course
	courseUUID, _ := CreateCourse("Set criteria to student grade test - course")

	// Get the invitation code of the course
	courseInvitationCode, _ := GetInvitationCode(courseUUID)

	// Create a laboratory
	laboratoryName := "Set criteria to student grade test - laboratory"
	laboratoryOpeningDate := "2023-12-01T08:00"
	laboratoryDueDate := "2023-12-01T12:00"

	laboratoryCreationResponse, _ := CreateLaboratory(cookie, map[string]interface{}{
		"name":         laboratoryName,
		"course_uuid":  courseUUID,
		"opening_date": laboratoryOpeningDate,
		"due_date":     laboratoryDueDate,
	})
	laboratoryUUID := laboratoryCreationResponse["uuid"].(string)

	// Create a rubric
	rubricName := "Set criteria to student grade test - rubric"
	rubricCreationResponse, _ := CreateRubric(cookie, map[string]interface{}{
		"name": rubricName,
	})
	rubricUUID := rubricCreationResponse["uuid"].(string)

	// Add an objective to the rubric
	objectiveName := "Set criteria to student grade test - objective"
	objectiveCreationResponse, _ := AddObjectiveToRubric(cookie, rubricUUID, map[string]interface{}{
		"description": objectiveName,
	})
	objectiveUUID := objectiveCreationResponse["uuid"].(string)

	// Add a criteria to the objective
	criteriaName := "Set criteria to student grade test - criteria"
	criteriaWeight := 1.0
	criteriaCreationResponse, _ := AddCriteriaToObjective(cookie, objectiveUUID, map[string]interface{}{
		"description": criteriaName,
		"weight":      criteriaWeight,
	})
	criteriaUUID := criteriaCreationResponse["uuid"].(string)

	// Add the rubric to the laboratory
	UpdateLaboratory(cookie, laboratoryUUID, map[string]interface{}{
		"rubric_uuid":  rubricUUID,
		"name":         laboratoryName,
		"opening_date": laboratoryOpeningDate,
		"due_date":     laboratoryDueDate,
	})

	// Add the student to the course
	AddStudentToCourse(courseInvitationCode)

	// Get the uuid of the student
	enrolledStudentsResponse, _ := GetStudentsEnrolledInCourse(cookie, courseUUID)
	enrolledStudents := enrolledStudentsResponse["students"].([]interface{})
	firstStudent := enrolledStudents[0].(map[string]interface{})
	studentUUID := firstStudent["uuid"].(string)

	// ## Test execution
	// Select the criteria for the student grade
	_, code := SetCriteriaToStudentGrade(&SetCriteriaToStudentGradeUtilsDTO{
		LaboratoryUUID: laboratoryUUID,
		RubricUUID:     rubricUUID,
		StudentUUID:    studentUUID,
		ObjectiveUUID:  objectiveUUID,
		CriteriaUUID:   criteriaUUID,
	}, cookie)
	c.Equal(http.StatusNoContent, code)

	// Get the student's grade
	studentGradeResponse, code := GetSummarizedGrades(laboratoryUUID, cookie)
	c.Equal(http.StatusOK, code)

	studentsGrades := studentGradeResponse["grades"].([]interface{})
	c.Equal(1, len(studentsGrades))

	studentGrade := studentsGrades[0].(map[string]interface{})
	c.Equal(studentUUID, studentGrade["student_uuid"].(string))
	c.Equal(criteriaWeight, studentGrade["grade"].(float64))
}
