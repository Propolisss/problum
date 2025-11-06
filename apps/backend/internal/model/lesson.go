package model

import "time"

/*
CREATE TABLE IF NOT EXISTS lessons (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    name TEXT,
    description TEXT,
    position INTEGER,
    content TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
*/

type Lesson struct {
	ID          int       `db:"id"`
	CourseID    int       `db:"course_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Position    int       `db:"position"`
	Content     string    `db:"content"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
