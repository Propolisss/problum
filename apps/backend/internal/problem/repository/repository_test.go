package repository

import (
	"context"
	"testing"
	"time"

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
		CREATE TABLE problems (
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			lesson_id INTEGER NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
			name TEXT,
			statement TEXT,
			difficulty TEXT,
			time_limit INTERVAL,
			memory_limit BIGINT,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);
		INSERT INTO courses (name, description, tags, status) VALUES ('Go', 'Go course', ARRAY['backend'], 'published');
		INSERT INTO lessons (course_id, name, description, position, content) VALUES (1, 'Intro', 'Start', 1, 'Content');
		INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
		VALUES (1, 'Two Sum', 'Solve', 'easy', INTERVAL '2 seconds', 268435456);
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
				if got.ID != 1 || got.Name != "Two Sum" || got.TimeLimit != 2*time.Second {
					t.Fatalf("unexpected problem: %+v", got)
				}
			},
		},
		{
			name: "list by lesson id",
			run: func(t *testing.T) {
				got, err := repo.ListByLessonID(ctx, 1)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(got) != 1 || got[0].LessonID != 1 {
					t.Fatalf("unexpected problems: %+v", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}
