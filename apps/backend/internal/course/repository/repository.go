package repository

import (
	"context"
	"errors"
	"fmt"

	"problum/internal/database"
	"problum/internal/model"

	"github.com/rs/zerolog/log"
)

var ErrNotFound = errors.New("course not found")

type Repository struct {
	db *database.DB
}

func New(db *database.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List(ctx context.Context) ([]*model.Course, error) {
	query := `
	SELECT
		id,
		name,
		description,
		tags,
		status,
		created_at,
		updated_at
	FROM courses
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to created courses list query")
		return nil, fmt.Errorf("failed to created courses list query: %w", err)
	}
	defer rows.Close()

	courses := make([]*model.Course, 0)
	for rows.Next() {
		course := &model.Course{}
		if err := rows.Scan(
			&course.ID,
			&course.Name,
			&course.Description,
			&course.Tags,
			&course.Status,
			&course.CreatedAt,
			&course.UpdatedAt,
		); err != nil {
			log.Error().Err(err).Msg("Failed to scan course")
			return nil, fmt.Errorf("failed to scan course: %w", err)
		}

		courses = append(courses, course)
	}
	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Failed to iterate courses")
		return nil, fmt.Errorf("failed to iterate courses: %w", err)
	}

	return courses, nil
}

func (r *Repository) Get(ctx context.Context, id int) (*model.Course, error) {
	query := `
	SELECT
		id,
		name,
		description,
		tags,
		status,
		created_at,
		updated_at
	FROM courses
	WHERE id = $1
	`

	course := &model.Course{}
	if err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&course.ID,
		&course.Name,
		&course.Description,
		&course.Tags,
		&course.Status,
		&course.CreatedAt,
		&course.UpdatedAt,
	); err != nil {
		log.Error().Err(err).Msg("Failed to get course")
		return nil, fmt.Errorf("failed to get course: %w", err)
	}

	return course, nil
}
