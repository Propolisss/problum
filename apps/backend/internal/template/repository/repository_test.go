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
		CREATE TABLE templates (
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
			language TEXT,
			code TEXT,
			metadata JSONB,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);
		INSERT INTO courses (name, description, tags, status) VALUES ('Go', 'Go course', ARRAY['backend'], 'published');
		INSERT INTO lessons (course_id, name, description, position, content) VALUES (1, 'Intro', 'Start', 1, 'Content');
		INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
		VALUES (1, 'Two Sum', 'Solve', 'easy', INTERVAL '2 seconds', 268435456);
		INSERT INTO templates (problem_id, language, code, metadata)
		VALUES (1, 'go', 'package main', '{"stdin": true}'), (1, 'python', 'print()', '{"stdin": true}');
	`)
	t.Cleanup(cleanup)

	repo := New(db)

	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "get by problem id and language",
			run: func(t *testing.T) {
				got, err := repo.GetByProblemIDAndLanguage(ctx, 1, "go")
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.ProblemID != 1 || got.Language != "go" || got.Code != "package main" {
					t.Fatalf("unexpected template: %+v", got)
				}
			},
		},
		{
			name: "get languages by problem id",
			run: func(t *testing.T) {
				got, err := repo.GetLanguagesByProblemID(ctx, 1)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(got) != 2 {
					t.Fatalf("got %d languages, want %d", len(got), 2)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}
