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

func (r *Repository) Enroll(ctx context.Context, courseID, userID int) error {
	query := `
	INSERT INTO enrollments(course_id, user_id)
	VALUES ($1, $2)
	`

	if _, err := r.db.Pool.Exec(ctx, query, courseID, userID); err != nil {
		log.Error().Err(err).Int("course_id", courseID).Int("user_id", userID).Msg("Failed to enroll course")
		return fmt.Errorf("failed to enroll course: %w", err)
	}

	return nil
}

func (r *Repository) Get(ctx context.Context, courseID, userID int) (*model.Enrollment, error) {
	query := `
	SELECT
		id,
		course_id,
		user_id,
		created_at,
		updated_at
	FROM enrollments
	WHERE course_id = $1 AND user_id = $2
	`

	enrollment := &model.Enrollment{}
	if err := r.db.Pool.QueryRow(ctx, query, courseID, userID).Scan(
		&enrollment.ID,
		&enrollment.CourseID,
		&enrollment.UserID,
		&enrollment.CreatedAt,
		&enrollment.UpdatedAt,
	); err != nil {
		log.Error().Err(err).Int("course_id", courseID).Int("user_id", userID).Msg("Failed to get enrollment")
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}

	return enrollment, nil
}

func (r *Repository) GetListByUserID(ctx context.Context, userID int) ([]*model.Enrollment, error) {
	query := `
	SELECT
		id,
		course_id,
		user_id,
		created_at,
		updated_at
	FROM enrollments
	WHERE user_id = $1
	`

	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create query row for list by user id")
		return nil, fmt.Errorf("failed to create query row for list by user id: %w", err)
	}
	defer rows.Close()

	enrollments := make([]*model.Enrollment, 0)
	for rows.Next() {
		enrollment := &model.Enrollment{}
		if err := rows.Scan(
			&enrollment.ID,
			&enrollment.CourseID,
			&enrollment.UserID,
			&enrollment.CreatedAt,
			&enrollment.UpdatedAt,
		); err != nil {
			log.Error().Err(err).Msg("Failed to scan enrollment")
			return nil, fmt.Errorf("failed to scan enrollment: %w", err)
		}

		enrollments = append(enrollments, enrollment)
	}
	if err := rows.Err(); err != nil {
		log.Error().Err(err).Msg("Failed to iterate enrollments list by user id")
		return nil, fmt.Errorf("failed to iterate enrollments list by user id: %w", err)
	}

	return enrollments, nil
}
