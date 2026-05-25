package repository

import (
	"context"
	"encoding/json"
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
		CREATE TABLE tests (
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
			tests JSONB,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);
		INSERT INTO courses (name, description, tags, status) VALUES ('Go', 'Go course', ARRAY['backend'], 'published');
		INSERT INTO lessons (course_id, name, description, position, content) VALUES (1, 'Intro', 'Start', 1, 'Content');
		INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
		VALUES (1, 'Two Sum', 'Solve', 'easy', INTERVAL '2 seconds', 268435456);
		INSERT INTO tests (problem_id, tests) VALUES (1, '[{"input":[1,2],"output":3}]');
	`)
	t.Cleanup(cleanup)

	repo := New(db)
	got, err := repo.GetByProblemID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ProblemID != 1 || !json.Valid(got.Tests) {
		t.Fatalf("unexpected test: %+v", got)
	}
}
