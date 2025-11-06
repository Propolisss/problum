package model

import "time"

/*
CREATE TABLE IF NOT EXISTS attempts (
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
*/

type Attempt struct {
	ID           int           `db:"id"`
	UserID       int           `db:"user_id"`
	ProblemID    int           `db:"problem_id"`
	Duration     time.Duration `db:"duration"`
	MemoryUsage  int64         `db:"memory_usage"`
	Language     string        `db:"language"`
	Code         string        `db:"code"`
	Status       string        `db:"status"`
	ErrorMessage *string       `db:"error_message"`
	CreatedAt    time.Time     `db:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at"`
}
