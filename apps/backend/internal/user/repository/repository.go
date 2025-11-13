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

var ErrNotFound = errors.New("user not found")

type Repository struct {
	db *database.DB
}

func New(db *database.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) FindByLogin(ctx context.Context, login string) (*model.User, error) {
	query := `
	SELECT 
		id,
		login,
		hashed_password,
		created_at,
		updated_at
	FROM users
	WHERE login = $1
	`

	user := &model.User{}
	if err := r.db.Pool.QueryRow(ctx, query, login).Scan(
		&user.ID,
		&user.Login,
		&user.HashedPassword,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msg("User not found")
			return nil, ErrNotFound
		}

		log.Error().Err(err).Msg("Failed to find user by login")
		return nil, fmt.Errorf("failed to find user by login: %w", err)
	}

	return user, nil
}

func (r *Repository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	query := `
	INSERT INTO	users(login, hashed_password)
	VALUES ($1, $2)
	RETURNING 
		id,
		login,
		hashed_password,
		created_at,
		updated_at 
	`

	u := &model.User{}
	if err := r.db.Pool.QueryRow(ctx, query, user.Login, user.HashedPassword).Scan(
		&u.ID,
		&u.Login,
		&u.HashedPassword,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		log.Error().Err(err).Msg("Failed to create user")
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return u, nil
}

func (r *Repository) Get(ctx context.Context, userID int) (*model.User, error) {
	query := `
	SELECT 
		id,
		login,
		hashed_password,
		created_at,
		updated_at
	FROM users
	WHERE id = $1
	`

	user := &model.User{}
	if err := r.db.Pool.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Login,
		&user.HashedPassword,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error().Err(err).Msg("User not found")
			return nil, ErrNotFound
		}

		log.Error().Err(err).Msg("Failed to find user by user_id")
		return nil, fmt.Errorf("failed to find user by user_id: %w", err)
	}

	return user, nil
}
