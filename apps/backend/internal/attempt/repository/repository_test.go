package repository

import (
	"context"
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
		CREATE TABLE attempts (
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
			duration INTERVAL,
			memory_usage BIGINT,
			language TEXT,
			code TEXT,
			status TEXT,
			error_message TEXT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW(),
			updated_at TIMESTAMPTZ DEFAULT NOW()
		);
		INSERT INTO users (login, hashed_password) VALUES ('user', 'hash');
		INSERT INTO courses (name, description, tags, status) VALUES ('Go', 'Go course', ARRAY['backend'], 'published');
		INSERT INTO lessons (course_id, name, description, position, content) VALUES (1, 'Intro', 'Start', 1, 'Content');
		INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
		VALUES (1, 'Two Sum', 'Solve', 'easy', INTERVAL '2 seconds', 268435456);
		INSERT INTO attempts (user_id, problem_id, duration, memory_usage, language, code, status)
		VALUES (1, 1, INTERVAL '1 second', 1024, 'go', 'package main', 'ok');
	`)
	t.Cleanup(cleanup)

	repo := New(db)

	tests := []struct {
		name string
		run  func(*testing.T)
	}{
		{
			name: "list by problem id",
			run: func(t *testing.T) {
				got, err := repo.ListByProblemID(ctx, 1, 1)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(got) != 1 || got[0].ProblemID != 1 {
					t.Fatalf("unexpected attempts: %+v", got)
				}
			},
		},
		{
			name: "list by user id",
			run: func(t *testing.T) {
				got, err := repo.ListByUserID(ctx, 1)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if len(got) != 1 || got[0].UserID != 1 {
					t.Fatalf("unexpected attempts: %+v", got)
				}
			},
		},
		{
			name: "submit",
			run: func(t *testing.T) {
				id, err := repo.Submit(ctx, &model.Attempt{
					UserID:      1,
					ProblemID:   1,
					Duration:    time.Second,
					MemoryUsage: 2048,
					Language:    "python",
					Code:        "print()",
					Status:      "pending",
				})
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if id == 0 {
					t.Fatal("expected id")
				}
			},
		},
		{
			name: "update",
			run: func(t *testing.T) {
				msg := "wrong answer"
				if err := repo.Update(ctx, &model.Attempt{
					ID:           1,
					Duration:     2 * time.Second,
					MemoryUsage:  4096,
					Language:     "go",
					Code:         "package main",
					Status:       "failed",
					ErrorMessage: &msg,
				}); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				got, err := repo.Get(ctx, 1)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if got.Status != "failed" || got.ErrorMessage == nil || *got.ErrorMessage != msg {
					t.Fatalf("unexpected attempt: %+v", got)
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
				if got.ID != 1 || got.UserID != 1 {
					t.Fatalf("unexpected attempt: %+v", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.run)
	}
}
