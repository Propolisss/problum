package repository

import (
	"context"
	"fmt"

	"problum/internal/database"
	"problum/internal/model"

	"github.com/rs/zerolog/log"
)

type Repository struct {
	db *database.DB
}

func New(db *database.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Get(ctx context.Context, id int) (*model.Problem, error) {
	query := `
	SELECT
		id,
		lesson_id,
		name,
		statement,
		difficulty,
		time_limit,
		memory_limit,
		created_at,
		updated_at
	FROM problems
	WHERE id = $1
	`

	problem := &model.Problem{}
	if err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&problem.ID,
		&problem.LessonID,
		&problem.Name,
		&problem.Statement,
		&problem.Difficulty,
		&problem.TimeLimit,
		&problem.MemoryLimit,
		&problem.CreatedAt,
		&problem.UpdatedAt,
	); err != nil {
		log.Error().Err(err).Msg("Failed to get problem")
		return nil, fmt.Errorf("failed to get problem: %w", err)
	}

	return problem, nil
}

func (r *Repository) ListByLessonID(ctx context.Context, id int) ([]*model.Problem, error) {
	query := `
	SELECT
		id,
		lesson_id,
		name,
		statement,
		difficulty,
		time_limit,
		memory_limit,
		created_at,
		updated_at
	FROM problems
	WHERE lesson_id = $1
	`

	rows, err := r.db.Pool.Query(ctx, query, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create problems list by lesson id query")
		return nil, fmt.Errorf("failed to create problems list by lesson id query: %w", err)
	}
	defer rows.Close()

	problems := make([]*model.Problem, 0)
	for rows.Next() {
		problem := &model.Problem{}
		if err := r.db.Pool.QueryRow(ctx, query, id).Scan(
			&problem.ID,
			&problem.LessonID,
			&problem.Name,
			&problem.Statement,
			&problem.Difficulty,
			&problem.TimeLimit,
			&problem.MemoryLimit,
			&problem.CreatedAt,
			&problem.UpdatedAt,
		); err != nil {
			log.Error().Err(err).Msg("Failed to scan problem")
			return nil, fmt.Errorf("failed to scan problem: %w", err)
		}

		problems = append(problems, problem)
	}
	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Failed to iterate problems by lesson id")
		return nil, fmt.Errorf("failed to iterate problems by lesson id: %w", err)
	}

	return problems, nil
}
