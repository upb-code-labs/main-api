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

	rubric = &entities.Rubric{}
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
	`, uuid)
	if err != nil {
		return nil, err
	}

	// Save the objectives into a map
	objectivesUUIDs := make([]string, 0)
	objectivesMap := make(map[string]*entities.RubricObjective)
	for rows.Next() {
		objective := &entities.RubricObjective{}
		err = rows.Scan(&objective.UUID, &objective.RubricUUID, &objective.Description)
		if err != nil {
			return nil, err
		}

		objectivesUUIDs = append(objectivesUUIDs, objective.UUID)
		objectivesMap[objective.UUID] = objective
	}

	// Get the objectives' criteria
	rows, err = repository.Connection.QueryContext(ctx, `
		SELECT id, objective_id, description, weight
		FROM criteria
		WHERE objective_id = ANY($1)
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

		objectivesMap[criteria.ObjectiveUUID].Criteria = append(objectivesMap[criteria.ObjectiveUUID].Criteria, *criteria)
	}

	// Append the objectives to the rubric
	for _, objective := range objectivesMap {
		rubric.Objectives = append(rubric.Objectives, *objective)
	}

	return rubric, nil
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
