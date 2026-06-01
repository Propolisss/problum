package repository

import (
	"context"
	"testing"

	"problum/internal/database"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestRepository(t *testing.T) {
	ctx := context.Background()
	repo, cleanup := setupRepository(t, ctx)
	t.Cleanup(cleanup)

	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "list",
			run: func(t *testing.T) {
				got, err := repo.List(ctx)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(got) != 2 {
					t.Fatalf("got %d courses, want %d", len(got), 2)
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
				if got.ID != 1 || got.Name != "Go" || got.Tags[0] != "backend" {
					t.Fatalf("unexpected course: %+v", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.run(t)
		})
	}
}

func setupRepository(t *testing.T, ctx context.Context) (*Repository, func()) {
	t.Helper()

	container, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("problum_test"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("get connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("connect postgres: %v", err)
	}

	_, err = pool.Exec(ctx, `
		CREATE TABLE courses (
			id INTEGER PRIMARY KEY,
			name TEXT,
			description TEXT,
			tags TEXT[],
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);
		INSERT INTO courses (id, name, description, tags, status)
		VALUES
			(1, 'Go', 'Go course', ARRAY['backend'], 'published'),
			(2, 'Python', 'Python course', ARRAY['backend', 'scripts'], 'draft');
	`)
	if err != nil {
		t.Fatalf("seed database: %v", err)
	}

	db := &database.DB{Pool: pool}
	return New(db), func() {
		pool.Close()
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("terminate postgres container: %v", err)
		}
	}
}
