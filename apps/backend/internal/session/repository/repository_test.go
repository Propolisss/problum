package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"problum/internal/model"
	"problum/internal/utils"
)

func TestRepository(t *testing.T) {
	ctx := context.Background()
	db, cleanup := utils.SetupPostgres(t, ctx, `
		CREATE TABLE users (
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			login TEXT UNIQUE NOT NULL,
			hashed_password TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);
		CREATE TABLE user_sessions (
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			refresh_hash TEXT NOT NULL,
			previous_refresh_hash TEXT DEFAULT '',
			expires_at TIMESTAMPTZ,
			revoked BOOLEAN DEFAULT false,
			last_activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);
		INSERT INTO users (login, hashed_password) VALUES ('user', 'hash');
		INSERT INTO user_sessions (user_id, refresh_hash, previous_refresh_hash, expires_at)
		VALUES (1, 'refresh', 'previous', NOW() + INTERVAL '1 day');
	`)
	t.Cleanup(cleanup)

	repo := New(db)

	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "get by refresh hash",
			run: func(t *testing.T) {
				got, err := repo.GetByRefreshHash(ctx, "refresh")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.UserID != 1 || got.RefreshHash != "refresh" {
					t.Fatalf("unexpected session: %+v", got)
				}
			},
		},
		{
			name: "get by previous refresh hash",
			run: func(t *testing.T) {
				got, err := repo.GetByPreviousRefreshHash(ctx, "previous")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.UserID != 1 || got.PreviousRefreshHash != "previous" {
					t.Fatalf("unexpected session: %+v", got)
				}
			},
		},
		{
			name: "create",
			run: func(t *testing.T) {
				got, err := repo.Create(ctx, &model.UserSession{
					UserID:      1,
					RefreshHash: "created",
					ExpiresAt:   time.Now().Add(time.Hour),
				})
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.ID == 0 || got.RefreshHash != "created" {
					t.Fatalf("unexpected session: %+v", got)
				}
			},
		},
		{
			name: "update",
			run: func(t *testing.T) {
				got, err := repo.Update(ctx, &model.UserSession{
					ID:                  1,
					UserID:              1,
					RefreshHash:         "updated",
					PreviousRefreshHash: "refresh",
					ExpiresAt:           time.Now().Add(time.Hour),
					Revoked:             true,
					LastActivityAt:      time.Now(),
				})
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.RefreshHash != "updated" || !got.Revoked {
					t.Fatalf("unexpected session: %+v", got)
				}
			},
		},
		{
			name: "logout all",
			run: func(t *testing.T) {
				if err := repo.LogoutAll(ctx, 1); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				got, err := repo.GetByRefreshHash(ctx, "updated")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if !got.Revoked {
					t.Fatalf("unexpected session: %+v", got)
				}
			},
		},
		{
			name: "not found",
			run: func(t *testing.T) {
				_, err := repo.GetByRefreshHash(ctx, "missing")
				if !errors.Is(err, ErrNotFound) {
					t.Fatalf("got %v, want %v", err, ErrNotFound)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}
