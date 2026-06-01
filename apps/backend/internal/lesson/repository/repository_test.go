package repository

import (
	"context"
	"testing"

	"problum/internal/utils"
)

func TestRepository(t *testing.T) {
	ctx := context.Background()
	db, cleanup := utils.SetupPostgres(t, ctx, `
		CREATE TABLE courses (
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			name TEXT,
			description TEXT,
			tags TEXT[],
			status TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);
		CREATE TABLE lessons (
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
			name TEXT,
			description TEXT,
			position INTEGER,
			content TEXT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);
		INSERT INTO courses (name, description, tags, status) VALUES ('Go', 'Go course', ARRAY['backend'], 'published');
		INSERT INTO lessons (course_id, name, description, position, content)
		VALUES (1, 'Intro', 'Start', 1, 'Content'), (1, 'Next', 'Continue', 2, 'More');
	`)
	t.Cleanup(cleanup)

	repo := New(db)

	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "get",
			run: func(t *testing.T) {
				got, err := repo.Get(ctx, 1)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.ID != 1 || got.Name != "Intro" {
					t.Fatalf("unexpected lesson: %+v", got)
				}
			},
		},
		{
			name: "list",
			run: func(t *testing.T) {
				got, err := repo.List(ctx)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(got) != 2 {
					t.Fatalf("got %d lessons, want %d", len(got), 2)
				}
			},
		},
		{
			name: "list by course id",
			run: func(t *testing.T) {
				got, err := repo.ListByCourseID(ctx, 1)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(got) != 2 || got[0].CourseID != 1 {
					t.Fatalf("unexpected lessons: %+v", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}
