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

func (r *Repository) GetByProblemID(ctx context.Context, problemID int) (*model.Test, error) {
	query := `
	SELECT
		id,
		problem_id,
		tests,
		created_at,
		updated_at
	FROM tests
	WHERE problem_id = $1
	`

	test := &model.Test{}
	if err := r.db.Pool.QueryRow(ctx, query, problemID).Scan(
		&test.ID,
		&test.ProblemID,
		&test.Tests,
		&test.CreatedAt,
		&test.UpdatedAt,
	); err != nil {
		log.Error().Err(err).Msg("Failed to get test by problem id")
		return nil, fmt.Errorf("failed to get test by problem id")
	}

	return test, nil
}
