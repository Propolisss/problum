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

func (r *Repository) Get(ctx context.Context, id int) (*model.Lesson, error) {
	query := `
	SELECT
		id,
		course_id,
		name,
		description,
		position,
		content,
		created_at,
		updated_at
	FROM lessons
	WHERE id = $1
	`

	lesson := &model.Lesson{}
	if err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&lesson.ID,
		&lesson.CourseID,
		&lesson.Name,
		&lesson.Description,
		&lesson.Position,
		&lesson.Content,
		&lesson.CreatedAt,
		&lesson.UpdatedAt,
	); err != nil {
		log.Error().Err(err).Msg("Failed to get lesson")
		return nil, fmt.Errorf("failed to get lesson")
	}

	return lesson, nil
}

func (r *Repository) List(ctx context.Context) ([]*model.Lesson, error) {
	query := `
	SELECT
		id,
		course_id,
		name,
		description,
		position,
		content,
		created_at,
		updated_at
	FROM lessons
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create lessons list query")
		return nil, fmt.Errorf("failed to create lessons list query: %w", err)
	}
	defer rows.Close()

	lessons := make([]*model.Lesson, 0)
	for rows.Next() {
		lesson := &model.Lesson{}
		if err := rows.Scan(
			&lesson.ID,
			&lesson.CourseID,
			&lesson.Name,
			&lesson.Description,
			&lesson.Position,
			&lesson.Content,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
		); err != nil {
			log.Error().Err(err).Msg("Failed to scan lesson")
			return nil, fmt.Errorf("failed to scan lesson: %w", err)
		}

		lessons = append(lessons, lesson)
	}
	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Failed to iterate lessons list")
		return nil, fmt.Errorf("failed to iterate lessons list: %w", err)
	}

	return lessons, nil
}

func (r *Repository) ListByCourseID(ctx context.Context, id int) ([]*model.Lesson, error) {
	query := `
	SELECT
		id,
		course_id,
		name,
		description,
		position,
		content,
		created_at,
		updated_at
	FROM lessons
	WHERE course_id = $1
	`

	rows, err := r.db.Pool.Query(ctx, query, id)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create lessons list by course id query")
		return nil, fmt.Errorf("failed to create lessons list course id query: %w", err)
	}
	defer rows.Close()

	lessons := make([]*model.Lesson, 0)
	for rows.Next() {
		lesson := &model.Lesson{}
		if err := rows.Scan(
			&lesson.ID,
			&lesson.CourseID,
			&lesson.Name,
			&lesson.Description,
			&lesson.Position,
			&lesson.Content,
			&lesson.CreatedAt,
			&lesson.UpdatedAt,
		); err != nil {
			log.Error().Err(err).Msg("Failed to scan lesson")
			return nil, fmt.Errorf("failed to scan lesson: %w", err)
		}

		lessons = append(lessons, lesson)
	}
	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Failed to iterate lessons list by course id")
		return nil, fmt.Errorf("failed to iterate lessons list course id: %w", err)
	}

	return lessons, nil
}
