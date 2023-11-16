package implementations

import (
	"context"
	"database/sql"
	"time"

	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/dtos"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/entities"
	"github.com/UPB-Code-Labs/main-api/src/rubrics/domain/errors"
	shared_infrastructure "github.com/UPB-Code-Labs/main-api/src/shared/infrastructure"
	"github.com/lib/pq"
)

type RubricsPostgresRepository struct {
	Connection *sql.DB
}

// Singleton
var rubricsPgRepositoryInstance *RubricsPostgresRepository

func GetRubricsPgRepository() *RubricsPostgresRepository {
	if rubricsPgRepositoryInstance == nil {
		rubricsPgRepositoryInstance = &RubricsPostgresRepository{
			Connection: shared_infrastructure.GetPostgresConnection(),
		}
	}

	return rubricsPgRepositoryInstance
}

func (repository *RubricsPostgresRepository) Save(dto *dtos.CreateRubricDTO) (rubric *entities.Rubric, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Start transaction
	tx, err := repository.Connection.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	// Create the rubric
	row := tx.QueryRowContext(ctx, `
			INSERT INTO rubrics (teacher_id, name)
			VALUES ($1, $2)
			RETURNING id
		`, dto.TeacherUUID, dto.Name)

	var rubricUUID string
	err = row.Scan(&rubricUUID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create an initial objective
	row = tx.QueryRowContext(ctx, `
			INSERT INTO objectives (rubric_id, description)
			VALUES ($1, $2)
			RETURNING id
		`, rubricUUID, "Initial objective")

	var objectiveUUID string
	err = row.Scan(&objectiveUUID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create an initial criteria
	_, err = tx.ExecContext(ctx, `
			INSERT INTO criteria (objective_id, description, weight)
			VALUES ($1, $2, $3)
		`, objectiveUUID, "Initial criteria", 5.00)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Return the rubric
	return repository.GetByUUID(rubricUUID)
}

func (repository *RubricsPostgresRepository) GetByUUID(uuid string) (rubric *entities.Rubric, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Get the rubric
	row := repository.Connection.QueryRowContext(ctx, `
		SELECT id, teacher_id, name
		FROM rubrics
		WHERE id = $1
	`, uuid)

	rubric = &entities.Rubric{
		Objectives: make([]entities.RubricObjective, 0),
	}

	err = row.Scan(&rubric.UUID, &rubric.TeacherUUID, &rubric.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &errors.RubricNotFoundError{}
		}

		return nil, err
	}

	// Get the rubric's objectives
	rows, err := repository.Connection.QueryContext(ctx, `
		SELECT id, rubric_id, description
		FROM objectives
		WHERE rubric_id = $1
		ORDER BY created_at ASC
	`, uuid)
	if err != nil {
		return nil, err
	}

	// Save the objectives into a map
	objectives := make([]*entities.RubricObjective, 0)
	objectivesUUIDs := make([]string, 0)
	objectivesIndex := make(map[string]int)

	for rows.Next() {
		objective := &entities.RubricObjective{
			Criteria: make([]entities.RubricObjectiveCriteria, 0),
		}

		err = rows.Scan(&objective.UUID, &objective.RubricUUID, &objective.Description)
		if err != nil {
			return nil, err
		}

		objectives = append(objectives, objective)
		objectivesUUIDs = append(objectivesUUIDs, objective.UUID)
		objectivesIndex[objective.UUID] = len(objectives) - 1
	}

	// Get the objectives' criteria
	rows, err = repository.Connection.QueryContext(ctx, `
		SELECT id, objective_id, description, weight
		FROM criteria
		WHERE objective_id = ANY($1)
		ORDER BY created_at ASC
	`, pq.Array(objectivesUUIDs))
	if err != nil {
		return nil, err
	}

	// Append the criteria to the objectives
	for rows.Next() {
		criteria := &entities.RubricObjectiveCriteria{}
		err = rows.Scan(&criteria.UUID, &criteria.ObjectiveUUID, &criteria.Description, &criteria.Weight)
		if err != nil {
			return nil, err
		}

		objectives[objectivesIndex[criteria.ObjectiveUUID]].Criteria = append(objectives[objectivesIndex[criteria.ObjectiveUUID]].Criteria, *criteria)
	}

	// Append the objectives to the rubric
	for _, objective := range objectives {
		rubric.Objectives = append(rubric.Objectives, *objective)
	}

	return rubric, nil
}

func (repository *RubricsPostgresRepository) DoesTeacherOwnRubric(teacherUUID string, rubricUUID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Get the rubric
	row := repository.Connection.QueryRowContext(ctx, `
		SELECT teacher_id
		FROM rubrics
		WHERE id = $1
	`, rubricUUID)

	var rubricTeacherUUID string
	err := row.Scan(&rubricTeacherUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, &errors.RubricNotFoundError{}
		}

		return false, err
	}

	return rubricTeacherUUID == teacherUUID, nil
}

func (repository *RubricsPostgresRepository) DoesTeacherOwnObjective(teacherUUID string, objectiveUUID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Get the objective
	row := repository.Connection.QueryRowContext(ctx, `
		SELECT teacher_id
		FROM objectives_owners
		WHERE objective_id = $1
	`, objectiveUUID)

	var objectiveTeacherUUID string
	err := row.Scan(&objectiveTeacherUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, &errors.ObjectiveNotFoundError{}
		}

		return false, err
	}

	return objectiveTeacherUUID == teacherUUID, nil
}

func (repository *RubricsPostgresRepository) DoesTeacherOwnCriteria(teacherUUID string, criteriaUUID string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Get the criteria
	row := repository.Connection.QueryRowContext(ctx, `
		SELECT teacher_id
		FROM criteria_owners
		WHERE criteria_id = $1
	`, criteriaUUID)

	var criteriaTeacherUUID string
	err := row.Scan(&criteriaTeacherUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, &errors.CriteriaNotFoundError{}
		}

		return false, err
	}

	return criteriaTeacherUUID == teacherUUID, nil
}

func (repository *RubricsPostgresRepository) GetAllCreatedByTeacher(teacherUUID string) ([]*dtos.CreatedRubricDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Get the rubrics
	rows, err := repository.Connection.QueryContext(ctx, `
		SELECT id, teacher_id, name
		FROM rubrics
		WHERE teacher_id = $1
	`, teacherUUID)
	if err != nil {
		return nil, err
	}

	rubrics := make([]*dtos.CreatedRubricDTO, 0)
	for rows.Next() {
		rubric := &dtos.CreatedRubricDTO{}
		err = rows.Scan(&rubric.UUID, &rubric.TeacherUUID, &rubric.Name)
		if err != nil {
			return nil, err
		}

		rubrics = append(rubrics, rubric)
	}

	return rubrics, nil
}

func (repository *RubricsPostgresRepository) AddObjectiveToRubric(rubricUUID string, objectiveDescription string) (objectiveUUID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Create the objective
	query := `
		INSERT INTO objectives (rubric_id, description)
		VALUES ($1, $2)
		RETURNING id
	`

	row := repository.Connection.QueryRowContext(ctx, query, rubricUUID, objectiveDescription)
	err = row.Scan(&objectiveUUID)
	if err != nil {
		return "", err
	}

	return objectiveUUID, nil
}

func (repository *RubricsPostgresRepository) UpdateObjective(dto *dtos.UpdateObjectiveDTO) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Update the objective
	query := `
		UPDATE objectives
		SET description = $1
		WHERE id = $2
	`

	_, err = repository.Connection.ExecContext(ctx, query, dto.UpdatedDescription, dto.ObjectiveUUID)
	if err != nil {
		return err
	}

	return nil
}

func (repository *RubricsPostgresRepository) AddCriteriaToObjective(dto *dtos.AddCriteriaToObjectiveDTO) (criteriaUUID string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Create the criteria
	query := `
		INSERT INTO criteria (objective_id, description, weight)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	row := repository.Connection.QueryRowContext(ctx, query, dto.ObjectiveUUID, dto.CriteriaDescription, dto.CriteriaWeight)
	err = row.Scan(&criteriaUUID)
	if err != nil {
		return "", err
	}

	return criteriaUUID, nil
}
