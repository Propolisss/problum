package repository

import (
	"context"
	"errors"
	"fmt"

	"problum/internal/database"
	"problum/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

var ErrNotFound = errors.New("session not found")

type Repository struct {
	db *database.DB
}

func New(db *database.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByRefreshHash(ctx context.Context, hash string) (*model.UserSession, error) {
	query := `
	SELECT
		id,
		user_id,
		refresh_hash,
		previous_refresh_hash,
		expires_at,
		revoked,
		last_activity_at,
		created_at,
		updated_at
	FROM user_sessions
	WHERE refresh_hash = $1
	`

	s := &model.UserSession{}
	if err := r.db.Pool.QueryRow(ctx, query, hash).Scan(
		&s.ID,
		&s.UserID,
		&s.RefreshHash,
		&s.PreviousRefreshHash,
		&s.ExpiresAt,
		&s.Revoked,
		&s.LastActivityAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msg("Session not found by refresh hash")
			return nil, ErrNotFound
		}

		log.Error().Err(err).Msg("Failed to get user session by refresh hash")
		return nil, fmt.Errorf("failed to get user session by refresh hash: %w", err)
	}

	return s, nil
}

func (r *Repository) GetByPreviousRefreshHash(ctx context.Context, hash string) (*model.UserSession, error) {
	query := `
	SELECT
		id,
		user_id,
		refresh_hash,
		previous_refresh_hash,
		expires_at,
		revoked,
		last_activity_at,
		created_at,
		updated_at
	FROM user_sessions
	WHERE previous_refresh_hash = $1
	`

	s := &model.UserSession{}
	if err := r.db.Pool.QueryRow(ctx, query, hash).Scan(
		&s.ID,
		&s.UserID,
		&s.RefreshHash,
		&s.PreviousRefreshHash,
		&s.ExpiresAt,
		&s.Revoked,
		&s.LastActivityAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msg("Session not found by previous refresh hash")
			return nil, ErrNotFound
		}

		log.Error().Err(err).Msg("Failed to get user session by previous refresh hash")
		return nil, fmt.Errorf("failed to get user session by previous refresh hash: %w", err)
	}

	return s, nil
}

func (r *Repository) Create(ctx context.Context, session *model.UserSession) (*model.UserSession, error) {
	query := `
	INSERT INTO user_sessions(
		user_id,
		refresh_hash,
		expires_at
	)
	VALUES (
		$1,
		$2,
		$3
	)
	RETURNING
		id,
		user_id,
		refresh_hash,
		previous_refresh_hash,
		expires_at,
		revoked,
		last_activity_at,
		created_at,
		updated_at
	`

	s := &model.UserSession{}
	if err := r.db.Pool.QueryRow(ctx, query, session.UserID, session.RefreshHash, session.ExpiresAt).Scan(
		&s.ID,
		&s.UserID,
		&s.RefreshHash,
		&s.PreviousRefreshHash,
		&s.ExpiresAt,
		&s.Revoked,
		&s.LastActivityAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	); err != nil {
		log.Error().Err(err).Msg("Failed to create user session")
		return nil, fmt.Errorf("failed to create user session: %w", err)
	}

	return s, nil
}

func (r *Repository) Update(ctx context.Context, session *model.UserSession) (*model.UserSession, error) {
	query := `
	UPDATE user_sessions
	SET
		user_id = $1,
		refresh_hash = $2,
		previous_refresh_hash = $3,
		expires_at = $4,
		revoked = $5,
		last_activity_at = $6
	WHERE id = $7
	RETURNING
		id,
		user_id,
		refresh_hash,
		previous_refresh_hash,
		expires_at,
		revoked,
		last_activity_at,
		created_at,
		updated_at
	`

	s := &model.UserSession{}
	if err := r.db.Pool.QueryRow(ctx, query,
		session.UserID,
		session.RefreshHash,
		session.PreviousRefreshHash,
		session.ExpiresAt,
		session.Revoked,
		session.LastActivityAt,
		session.ID,
	).Scan(
		&s.ID,
		&s.UserID,
		&s.RefreshHash,
		&s.PreviousRefreshHash,
		&s.ExpiresAt,
		&s.Revoked,
		&s.LastActivityAt,
		&s.CreatedAt,
		&s.UpdatedAt,
	); err != nil {
		log.Error().Err(err).Msg("Failed to update user session")
		return nil, fmt.Errorf("failed to update user session: %w", err)
	}

	return s, nil
}

func (r *Repository) LogoutAll(ctx context.Context, id int) error {
	query := `
	UPDATE user_sessions
	SET revoked = true
	WHERE user_id = $1	
	`

	if _, err := r.db.Pool.Exec(ctx, query, id); err != nil {
		log.Error().Err(err).Msg("Failed to logout all sessions")
		return fmt.Errorf("failed to logout all sessions: %w", err)
	}

	return nil
}
