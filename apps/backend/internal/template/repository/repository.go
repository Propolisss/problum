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

func (r *Repository) GetByProblemIDAndLanguage(
	ctx context.Context,
	problemID int,
	language string,
) (*model.Template, error) {
	query := `
	SELECT
		id,
		problem_id,
		language,
		code,
		metadata,
		created_at,
		updated_at
	FROM templates
	WHERE problem_id = $1 AND language = $2
	`

	template := &model.Template{}
	if err := r.db.Pool.QueryRow(ctx, query, problemID, language).Scan(
		&template.ID,
		&template.ProblemID,
		&template.Language,
		&template.Code,
		&template.Metadata,
		&template.CreatedAt,
		&template.UpdatedAt,
	); err != nil {
		log.Error().Err(err).Msg("Failed to get template by problem id and language")
		return nil, fmt.Errorf("failed to get template by problem id and language: %w", err)
	}

	return template, nil
}

func (r *Repository) GetLanguagesByProblemID(ctx context.Context, problemID int) ([]string, error) {
	query := `
		SELECT
			language
		FROM templates
		WHERE problem_id = $1	
	`

	rows, err := r.db.Pool.Query(ctx, query, problemID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create get languages query")
		return nil, fmt.Errorf("failed to create get languages query: %w", err)
	}
	defer rows.Close()

	languages := make([]string, 0)

	for rows.Next() {
		var language string

		if err := rows.Scan(&language); err != nil {
			log.Error().Err(err).Msg("Failed to scan language")
			return nil, fmt.Errorf("failed to scan language: %w", err)
		}

		languages = append(languages, language)
	}
	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Failed to iterate languages")
		return nil, fmt.Errorf("failed to iterate languages: %w", err)
	}

	return languages, nil
}
