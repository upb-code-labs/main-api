package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGradeStudentSubmission(t *testing.T) {
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

	// ## Test: Set criteria to student grade
	// Select the criteria for the student grade
	_, code := SetCriteriaToStudentGrade(&SetCriteriaToStudentGradeUtilsDTO{
		LaboratoryUUID: laboratoryUUID,
		StudentUUID:    studentUUID,
		ObjectiveUUID:  objectiveUUID,
		CriteriaUUID:   criteriaUUID,
	}, cookie)
	c.Equal(http.StatusNoContent, code)

	// ## Test: Set comment to student grade
	comment := "Set criteria to student grade test - comment"
	_, code = SetCommentToStudentGrade(&SetCommentToStudentGradeUtilsDTO{
		LaboratoryUUID: laboratoryUUID,
		StudentUUID:    studentUUID,
		Comment:        comment,
	}, cookie)
	c.Equal(http.StatusNoContent, code)

	// ## Test: Get all the grades of students in the laboratory
	studentGradeResponse, code := GetSummarizedGrades(laboratoryUUID, cookie)
	c.Equal(http.StatusOK, code)

	studentsGrades := studentGradeResponse["grades"].([]interface{})
	c.Equal(1, len(studentsGrades))

	firstStudentGrade := studentsGrades[0].(map[string]interface{})
	c.Equal(studentUUID, firstStudentGrade["student_uuid"].(string))
	c.Equal(criteriaWeight, firstStudentGrade["grade"].(float64))

	// ## Test: Get the grade of the student in the laboratory
	studentGradeResponse, code = GetStudentGrade(&GetStudentGradeUtilsDTO{
		LaboratoryUUID: laboratoryUUID,
		RubricUUID:     rubricUUID,
		StudentUUID:    studentUUID,
	}, cookie)
	c.Equal(http.StatusOK, code)

	studentGrade := studentGradeResponse["grade"].(float64)
	c.Equal(criteriaWeight, studentGrade)

	gradeComment := studentGradeResponse["comment"].(string)
	c.Equal(comment, gradeComment)

	selectedCriteriaList := studentGradeResponse["selected_criteria"].([]interface{})
	c.Equal(1, len(selectedCriteriaList))

	selectedCriteria := selectedCriteriaList[0].(map[string]interface{})
	selectedCriteriaObjectiveUUID := selectedCriteria["objective_uuid"].(string)
	c.Equal(objectiveUUID, selectedCriteriaObjectiveUUID)
	selectedCriteriaCriteriaUUID := selectedCriteria["criteria_uuid"].(string)
	c.Equal(criteriaUUID, selectedCriteriaCriteriaUUID)
}
