-- +goose Up
-- +goose StatementBegin
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
