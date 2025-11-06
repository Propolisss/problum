package model

import "time"

/*
CREATE TABLE IF NOT EXISTS problems (
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
*/

type Problem struct {
	ID          int           `db:"id"`
	LessonID    int           `db:"lesson_id"`
	Name        string        `db:"name"`
	Statement   string        `db:"statement"`
	Difficulty  string        `db:"difficulty"`
	TimeLimit   time.Duration `db:"time_limit"`
	MemoryLimit int64         `db:"memory_limit"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
}
