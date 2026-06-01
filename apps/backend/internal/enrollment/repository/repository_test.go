package repository

import (
	"context"
	"testing"

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
		CREATE TABLE courses (
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			name TEXT,
			description TEXT,
			tags TEXT[],
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);
		CREATE TABLE enrollments (
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW(),
			UNIQUE(course_id, user_id)
		);
		INSERT INTO users (login, hashed_password) VALUES ('user', 'hash');
		INSERT INTO courses (name, description, tags, status) VALUES ('Go', 'Go course', ARRAY['backend'], 'published');
	`)
	t.Cleanup(cleanup)

	repo := New(db)

	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "enroll",
			run: func(t *testing.T) {
				if err := repo.Enroll(ctx, 1, 1); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
		{
			name: "get",
			run: func(t *testing.T) {
				got, err := repo.Get(ctx, 1, 1)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.CourseID != 1 || got.UserID != 1 {
					t.Fatalf("unexpected enrollment: %+v", got)
				}
			},
		},
		{
			name: "list by user id",
			run: func(t *testing.T) {
				got, err := repo.GetListByUserID(ctx, 1)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(got) != 1 || got[0].CourseID != 1 {
					t.Fatalf("unexpected enrollments: %+v", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}
