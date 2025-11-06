package model

import (
	"encoding/json"
	"time"
)

/*
CREATE TABLE IF NOT EXISTS templates (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    language TEXT,
    code TEXT,
	metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
*/

type Template struct {
	ID        int             `db:"id"`
	ProblemID int             `db:"problem_id"`
	Language  string          `db:"language"`
	Code      string          `db:"code"`
	Metadata  json.RawMessage `db:"metadata"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}
