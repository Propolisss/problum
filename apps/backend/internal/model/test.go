package model

import (
	"encoding/json"
	"time"
)

/*
CREATE TABLE IF NOT EXISTS tests (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    tests JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
)
*/

type Test struct {
	ID        int             `db:"id"`
	ProblemID int             `db:"problem_id"`
	Tests     json.RawMessage `db:"tests"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}
