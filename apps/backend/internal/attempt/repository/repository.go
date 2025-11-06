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

func (r *Repository) ListByProblemID(ctx context.Context, userID, problemID int) ([]*model.Attempt, error) {
	query := `
	SELECT
		id,
		user_id,
		problem_id,
		duration,
		memory_usage,
		language,
		code,
		status,
		error_message,
		created_at,
		updated_at
	FROM attempts
	WHERE user_id = $1 AND problem_id = $2
	`

	rows, err := r.db.Pool.Query(ctx, query, userID, problemID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create attempts list by problem id query")
		return nil, fmt.Errorf("failed to create attempts list by problem id query: %w", err)
	}
	defer rows.Close()

	attempts := make([]*model.Attempt, 0)

	for rows.Next() {
		attempt := &model.Attempt{}
		if err := rows.Scan(
			&attempt.ID,
			&attempt.UserID,
			&attempt.ProblemID,
			&attempt.Duration,
			&attempt.MemoryUsage,
			&attempt.Language,
			&attempt.Code,
			&attempt.Status,
			&attempt.ErrorMessage,
			&attempt.CreatedAt,
			&attempt.UpdatedAt,
		); err != nil {
			log.Error().Err(err).Msg("Failed to scan attempt")
			return nil, fmt.Errorf("failed to scan attempt: %w", err)
		}

		attempts = append(attempts, attempt)
	}
	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Failed to iterate attempts by problem id")
		return nil, fmt.Errorf("failed to iterate attempts by problem id: %w", err)
	}

	return attempts, nil
}

func (r *Repository) ListByUserID(ctx context.Context, userID int) ([]*model.Attempt, error) {
	query := `
	SELECT
		id,
		user_id,
		problem_id,
		duration,
		memory_usage,
		language,
		code,
		status,
		error_message,
		created_at,
		updated_at
	FROM attempts
	WHERE user_id = $1
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create attempts list by user id query")
		return nil, fmt.Errorf("failed to create attempts list by user id query: %w", err)
	}
	defer rows.Close()

	attempts := make([]*model.Attempt, 0)

	for rows.Next() {
		attempt := &model.Attempt{}
		if err := rows.Scan(
			&attempt.ID,
			&attempt.UserID,
			&attempt.ProblemID,
			&attempt.Duration,
			&attempt.MemoryUsage,
			&attempt.Language,
			&attempt.Code,
			&attempt.Status,
			&attempt.ErrorMessage,
			&attempt.CreatedAt,
			&attempt.UpdatedAt,
		); err != nil {
			log.Error().Err(err).Msg("Failed to scan attempt")
			return nil, fmt.Errorf("failed to scan attempt: %w", err)
		}

		attempts = append(attempts, attempt)
	}
	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Failed to iterate attempts by user id")
		return nil, fmt.Errorf("failed to iterate attempts by user id: %w", err)
	}

	return attempts, nil
}

func (r *Repository) Submit(ctx context.Context, attempt *model.Attempt) (int, error) {
	query := `
	INSERT INTO attempts(
		user_id,
		problem_id,
		duration,
		memory_usage,
		language,
		code,
		status,
		error_message
	)
	VALUES(
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8
	)
	RETURNING
		id
	`

	id := 0
	if err := r.db.Pool.QueryRow(ctx, query,
		attempt.UserID,
		attempt.ProblemID,
		attempt.Duration,
		attempt.MemoryUsage,
		attempt.Language,
		attempt.Code,
		attempt.Status,
		attempt.ErrorMessage,
	).Scan(&id); err != nil {
		log.Error().Err(err).Msg("Failed to insert attempt")
		return 0, fmt.Errorf("failed to insert attempt: %w", err)
	}

	return id, nil
}

func (r *Repository) Update(ctx context.Context, attempt *model.Attempt) error {
	query := `
	UPDATE attempts
	SET
		duration = $1,
		memory_usage = $2,
		language = $3,
		code = $4,
		status = $5,
		error_message = $6
	WHERE id = $7
	`

	if _, err := r.db.Pool.Exec(ctx, query,
		attempt.Duration,
		attempt.MemoryUsage,
		attempt.Language,
		attempt.Code,
		attempt.Status,
		attempt.ErrorMessage,
		attempt.ID,
	); err != nil {
		log.Error().Err(err).Msg("Failed to update attempt")
		return fmt.Errorf("failed to update attempt: %w", err)
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, attemptID int) (*model.Attempt, error) {
	query := `
	SELECT
		id,
		user_id,
		problem_id,
		duration,
		memory_usage,
		language,
		code,
		status,
		error_message,
		created_at,
		updated_at
	FROM attempts
	WHERE id = $1
	`

	attempt := &model.Attempt{}
	if err := r.db.Pool.QueryRow(ctx, query, attemptID).Scan(
		&attempt.ID,
		&attempt.UserID,
		&attempt.ProblemID,
		&attempt.Duration,
		&attempt.MemoryUsage,
		&attempt.Language,
		&attempt.Code,
		&attempt.Status,
		&attempt.ErrorMessage,
		&attempt.CreatedAt,
		&attempt.UpdatedAt,
	); err != nil {
		log.Error().Err(err).Msg("Failed to get attempt")
		return nil, fmt.Errorf("failed to get attempt: %w", err)
	}

	return attempt, nil
}
