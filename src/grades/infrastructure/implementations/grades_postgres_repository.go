package implementations

import (
	"context"
	"database/sql"
	"time"

	"github.com/UPB-Code-Labs/main-api/src/grades/domain/dtos"
	sharedInfrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
)

// GradesPostgresRepository implementation of the GradesRepository interface
type GradesPostgresRepository struct {
	Connection *sql.DB
}

var gradesRepositoryInstance *GradesPostgresRepository

// GetGradesPostgresRepositoryInstance returns the singleton instance of the GradesPostgresRepository
func GetGradesPostgresRepositoryInstance() *GradesPostgresRepository {
	if gradesRepositoryInstance == nil {
		gradesRepositoryInstance = &GradesPostgresRepository{
			Connection: sharedInfrastructure.GetPostgresConnection(),
		}
	}

	return gradesRepositoryInstance
}

// GetStudentsGradesInLaboratory returns the grades of the students in a laboratory
// that were graded using the current rubric of the laboratory by the teacher
func (repository *GradesPostgresRepository) GetStudentsGradesInLaboratory(laboratoryUUID, rubricUUID string) ([]*dtos.SummarizedStudentGradeDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	query := `
		SELECT student_id, student_full_name, total_criteria_weight
		FROM summarized_grades
		WHERE laboratory_id = $1 AND rubric_id = $2
	`

	// Run the query
	rows, err := repository.Connection.QueryContext(ctx, query, laboratoryUUID, rubricUUID)
	if err != nil {
		return nil, err
	}

	// Parse the results
	var summarizedGrades []*dtos.SummarizedStudentGradeDTO
	for rows.Next() {
		var studentGrade dtos.SummarizedStudentGradeDTO

		if err := rows.Scan(
			&studentGrade.StudentUUID,
			&studentGrade.StudentFullName,
			&studentGrade.Grade); err != nil {
			return nil, err
		}

		summarizedGrades = append(summarizedGrades, &studentGrade)
	}

	return summarizedGrades, nil
}

// SetCriteriaToGrade sets a criteria to a student's grade
func (repository *GradesPostgresRepository) SetCriteriaToGrade(dto *dtos.SetCriteriaToGradeDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Check if the student has a grade in the laboratory
	studentHasGrade, err := repository.doesStudentHasGrade(&dtos.CheckIfStudentHasGradeDTO{
		StudentUUID:    dto.StudentUUID,
		LaboratoryUUID: dto.LaboratoryUUID,
		RubricUUID:     dto.RubricUUID,
	})
	if err != nil {
		return err
	}

	// Get the UUID of the grade of the student in the laboratory with the given rubric
	studentGradeUUID := ""

	if !studentHasGrade {
		// Create a grade for the student if they do not have one
		studentGradeUUID, err = repository.createStudentGrade(&dtos.CreateStudentGradeDTO{
			CheckIfStudentHasGradeDTO: dtos.CheckIfStudentHasGradeDTO{
				StudentUUID:    dto.StudentUUID,
				LaboratoryUUID: dto.LaboratoryUUID,
				RubricUUID:     dto.RubricUUID,
			},
		})

		if err != nil {
			return err
		}
	} else {
		// Get the UUID of the grade of the student in the laboratory with the given rubric
		studentGradeUUID, err = repository.getStudentGradeUUID(&dtos.GetStudentGradeDTO{
			CheckIfStudentHasGradeDTO: dtos.CheckIfStudentHasGradeDTO{
				StudentUUID:    dto.StudentUUID,
				LaboratoryUUID: dto.LaboratoryUUID,
				RubricUUID:     dto.RubricUUID,
			},
		})

		if err != nil {
			return err
		}
	}

	// UPSERT the criteria to the grade
	query := `
		INSERT INTO grade_has_criteria (grade_id, criteria_id, objective_id)	
		VALUES ($1, $2, $3)
		ON CONFLICT (grade_id, objective_id) DO
		UPDATE SET 
			criteria_id = $2
	`

	// Run the query
	if _, err := repository.Connection.ExecContext(
		ctx,
		query,
		studentGradeUUID,
		dto.CriteriaUUID,
		dto.ObjectiveUUID,
	); err != nil {
		return err
	}

	return nil
}

// doesStudentHasGrade checks if a student has a grade in a laboratory
func (repository *GradesPostgresRepository) doesStudentHasGrade(dto *dtos.CheckIfStudentHasGradeDTO) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM grades
			WHERE student_id = $1 AND laboratory_id = $2 AND rubric_id = $3
		)
	`

	// Run the query
	row := repository.Connection.QueryRowContext(
		ctx,
		query,
		dto.StudentUUID,
		dto.LaboratoryUUID,
		dto.RubricUUID,
	)

	// Parse the result
	var studentHasGrade bool
	if err := row.Scan(&studentHasGrade); err != nil {
		return false, err
	}

	return studentHasGrade, nil
}

// getStudentGrade returns the grade of a student in a laboratory
func (repository *GradesPostgresRepository) getStudentGradeUUID(dto *dtos.GetStudentGradeDTO) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	query := `
		SELECT id
		FROM grades
		WHERE student_id = $1 AND laboratory_id = $2 AND rubric_id = $3
	`

	// Run the query
	row := repository.Connection.QueryRowContext(
		ctx,
		query,
		dto.StudentUUID,
		dto.LaboratoryUUID,
		dto.RubricUUID,
	)

	// Parse the result
	var gradeUUID string
	if err := row.Scan(&gradeUUID); err != nil {
		return "", err
	}

	return gradeUUID, nil
}

// createStudentGrade creates a grade for a student in a laboratory
func (repository *GradesPostgresRepository) createStudentGrade(dto *dtos.CreateStudentGradeDTO) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	query := `
		INSERT INTO grades (student_id, laboratory_id, rubric_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	// Run the query
	row := repository.Connection.QueryRowContext(
		ctx,
		query,
		dto.StudentUUID,
		dto.LaboratoryUUID,
		dto.RubricUUID,
	)

	// Parse the result
	var gradeUUID string
	if err := row.Scan(&gradeUUID); err != nil {
		return "", err
	}

	return gradeUUID, nil
}

// GetStudentGradeInLaboratoryWithRubric returns the grade of an student in a laboratory
// that was graded with an specific rubric
func (repository *GradesPostgresRepository) GetStudentGradeInLaboratoryWithRubric(dto *dtos.GetStudentGradeInLaboratoryWithRubricDTO) (*dtos.StudentGradeInLaboratoryWithRubricDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Get the grade and comment
	query := `
		SELECT grade_id, total_criteria_weight, comment
		FROM summarized_grades
		WHERE student_id = $1 AND laboratory_id = $2 AND rubric_id = $3
	`

	row := repository.Connection.QueryRowContext(
		ctx,
		query,
		dto.StudentUUID,
		dto.LaboratoryUUID,
		dto.RubricUUID,
	)

	var gradeUUID string
	grade := &dtos.StudentGradeInLaboratoryWithRubricDTO{}

	if err := row.Scan(&gradeUUID, &grade.Grade, &grade.Comment); err != nil {
		// If the student does not have a grade, return the zero-valued grade
		if err == sql.ErrNoRows {
			grade.SelectedCriteria = []*dtos.SelectedCriteriaInStudentGradeDTO{}
			return grade, nil
		}

		return nil, err
	}

	// Get the selected criteria
	query = `
		SELECT criteria_id, objective_id
		FROM grade_has_criteria
		WHERE grade_id = $1
	`

	rows, err := repository.Connection.QueryContext(ctx, query, gradeUUID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var selectedCriteria dtos.SelectedCriteriaInStudentGradeDTO

		if err := rows.Scan(&selectedCriteria.CriteriaUUID, &selectedCriteria.ObjectiveUUID); err != nil {
			return nil, err
		}

		grade.SelectedCriteria = append(grade.SelectedCriteria, &selectedCriteria)
	}

	return grade, nil
}

// SetCommentToGrade sets a comment to a student's grade
func (repository *GradesPostgresRepository) SetCommentToGrade(dto *dtos.SetCommentToGradeDTO) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// Check if the student has a grade in the laboratory
	studentHasGrade, err := repository.doesStudentHasGrade(&dtos.CheckIfStudentHasGradeDTO{
		StudentUUID:    dto.StudentUUID,
		LaboratoryUUID: dto.LaboratoryUUID,
		RubricUUID:     dto.RubricUUID,
	})
	if err != nil {
		return err
	}

	// Get the UUID of the grade of the student in the laboratory with the given rubric
	studentGradeUUID := ""

	if !studentHasGrade {
		// Create a grade for the student if they do not have one
		studentGradeUUID, err = repository.createStudentGrade(&dtos.CreateStudentGradeDTO{
			CheckIfStudentHasGradeDTO: dtos.CheckIfStudentHasGradeDTO{
				StudentUUID:    dto.StudentUUID,
				LaboratoryUUID: dto.LaboratoryUUID,
				RubricUUID:     dto.RubricUUID,
			},
		})

		if err != nil {
			return err
		}
	} else {
		// Get the UUID of the grade of the student in the laboratory with the given rubric
		studentGradeUUID, err = repository.getStudentGradeUUID(&dtos.GetStudentGradeDTO{
			CheckIfStudentHasGradeDTO: dtos.CheckIfStudentHasGradeDTO{
				StudentUUID:    dto.StudentUUID,
				LaboratoryUUID: dto.LaboratoryUUID,
				RubricUUID:     dto.RubricUUID,
			},
		})

		if err != nil {
			return err
		}
	}

	// Set the comment to the grade
	query := `
		UPDATE grades
		SET comment = $1
		WHERE id = $2
	`

	// Run the query
	if _, err := repository.Connection.ExecContext(
		ctx,
		query,
		dto.Comment,
		studentGradeUUID,
	); err != nil {
		return err
	}

	return nil
}
