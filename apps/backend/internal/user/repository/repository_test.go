package repository

import (
	"context"
	"errors"
	"testing"

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
		INSERT INTO users (login, hashed_password) VALUES ('user', 'hash');
	`)
	t.Cleanup(cleanup)

	repo := New(db)

	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "find by login",
			run: func(t *testing.T) {
				got, err := repo.FindByLogin(ctx, "user")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.Login != "user" || got.HashedPassword != "hash" {
					t.Fatalf("unexpected user: %+v", got)
				}
			},
		},
		{
			name: "create",
			run: func(t *testing.T) {
				got, err := repo.Create(ctx, &model.User{Login: "new", HashedPassword: "new-hash"})
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.ID == 0 || got.Login != "new" || got.HashedPassword != "new-hash" {
					t.Fatalf("unexpected user: %+v", got)
				}
			},
		},
		{
			name: "get",
			run: func(t *testing.T) {
				got, err := repo.Get(ctx, 1)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.ID != 1 || got.Login != "user" {
					t.Fatalf("unexpected user: %+v", got)
				}
			},
		},
		{
			name: "not found",
			run: func(t *testing.T) {
				_, err := repo.Get(ctx, 100)
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
