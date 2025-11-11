-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    login TEXT UNIQUE NOT NULL CHECK (
        length (login) > 0
        AND length (login) <= 50
    ),
    hashed_password TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_sessions (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_hash TEXT NOT NULL,
    previous_refresh_hash TEXT DEFAULT '',
    expires_at TIMESTAMPTZ,
    -- device_info TEXT,
    -- last_ip INET,
    revoked BOOLEAN DEFAULT false,
    last_activity_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS courses (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT,
    description TEXT,
    tags TEXT [],
    status TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

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

CREATE TABLE IF NOT EXISTS enrollments (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    course_id INTEGER NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(course_id, user_id)
);

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

CREATE TABLE IF NOT EXISTS tests (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    tests JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

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

CREATE TABLE IF NOT EXISTS templates (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    problem_id INTEGER NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
    language TEXT,
    code TEXT,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_sessions_user_id ON user_sessions(user_id);

CREATE INDEX IF NOT EXISTS idx_user_sessions_refresh_hash ON user_sessions(refresh_hash);

CREATE INDEX IF NOT EXISTS idx_lessons_course_id_position ON lessons(course_id, position);

CREATE INDEX IF NOT EXISTS idx_enrollments_user_id_course_id ON enrollments(user_id, course_id);

CREATE INDEX IF NOT EXISTS idx_problems_lesson_id ON problems(lesson_id);

CREATE INDEX IF NOT EXISTS idx_attempts_user_id_problem_id ON attempts(user_id, problem_id);

CREATE INDEX IF NOT EXISTS idx_templates_problem_id_language ON templates(problem_id, language);

CREATE INDEX IF NOT EXISTS idx_courses_tags_gin ON courses USING GIN(tags);

BEGIN;

WITH course_insert AS (
    INSERT INTO courses (name, description, tags, status)
    VALUES ('Основы Go', 'Базовый курс для изучения основ языка программирования Go.', '{"go", "beginner"}', 'published')
    RETURNING id
),
lesson_insert AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES ((SELECT id FROM course_insert), 'Первая программа', 'Знакомство с базовыми операциями и структурой программы на Go.', 1, '<h1>Введение</h1><p>В этом уроке мы напишем простую функцию.</p>')
    RETURNING id
),
problem_insert AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES ((SELECT id FROM lesson_insert), 'Сумма двух чисел', 'Напишите функцию `sum`, которая принимает два целых числа `a` и `b` и возвращает их сумму.', 'easy', '1 second', 67108864)
    RETURNING id
),
template_insert_go AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES ((SELECT id FROM problem_insert), 'go',
    '
// sum принимает два целых числа и должна вернуть их сумму.
func sum(a, b int) int {
    // Вставьте ваш код здесь
}',
    '{
      "function_name": "sum",
      "parameters": [
        { "name": "a", "type": "int" },
        { "name": "b", "type": "int" }
      ]
    }'::jsonb)
),
template_insert_python AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES ((SELECT id FROM problem_insert), 'python',
    '
# sum принимает два числа и должна вернуть их сумму.
def sum(a, b):
    # Вставьте ваш код здесь
',
    '{
      "function_name": "sum",
      "parameters": [
        { "name": "a", "type": "int" },
        { "name": "b", "type": "int" }
      ]
    }'::jsonb)
)
INSERT INTO tests (problem_id, tests)
VALUES ((SELECT id FROM problem_insert),
'[
  { "input": { "a": 1, "b": 2 }, "output": 3 },
  { "input": { "a": -5, "b": 5 }, "output": 0 },
  { "input": { "a": 1000000, "b": 2000000 }, "output": 3000000 },
  { "input": { "a": 0, "b": 0 }, "output": 0 },
  { "input": { "a": -10, "b": -20 }, "output": -30 }
]'::jsonb);

COMMIT;

-- +goose StatementEnd
