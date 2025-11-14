-- +goose Up
-- +goose StatementBegin
BEGIN;

WITH course_insert AS (
    INSERT INTO courses (name, description, tags, status)
    VALUES (
        'Основы Python для начинающих: От синтаксиса до структур данных',
        'Этот курс предназначен для тех, кто только начинает свой путь в программировании. Мы рассмотрим основы синтаксиса Python, структуры данных и базовые концепции, необходимые для решения практических задач.',
        '{"python", "beginner", "programming", "core"}',
        'published'
    )
    RETURNING id
),

lesson_1 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Введение в Python, переменные и операции',
        'Знакомство с базовым синтаксисом, выводом данных, переменными, основными типами и арифметическими операциями.',
        1,
        '
        <h1>Урок 1: Введение в Python, переменные и операции</h1>
        <p>Python — это высокоуровневый язык программирования, который ценится за свою <b>читаемость</b> и <b>простоту</b>. Он позволяет разработчикам выражать концепции меньшим количеством строк кода по сравнению с другими языками.</p>

        <h2>Комментарии</h2>
        <p>Комментарии — это строки, которые игнорируются интерпретатором, но очень важны для разработчиков, так как они объясняют логику кода. Однострочные комментарии начинаются с символа решетки (<code>#</code>).</p>
        <pre><code># Это однострочный комментарий
print("Начало программы") # Комментарий в конце строки</code></pre>

        <h2>Вывод данных: Функция <code>print()</code></h2>
        <p>Для отображения информации в консоли используется встроенная функция <code>print()</code>. Она может выводить строки, числа, переменные и даже несколько объектов сразу, разделяя их пробелом.</p>
        <pre><code>print("Hello, World!")
name = "Алиса"
print("Привет,", name, "!") # Выведет: Привет, Алиса !</code></pre>

        <h2>Переменные и присваивание</h2>
        <p>Переменная — это именованное место в памяти для хранения данных. В Python переменные создаются динамически при первом присваивании с помощью оператора <code>=</code>.</p>
        <pre><code>user_age = 25
gpa = 4.0
is_student = True</code></pre>

        <h2>Основные типы данных</h2>
        <p>Python является <b>динамически типизированным</b> языком, но сами данные имеют строгие типы. Вы можете проверить тип любой переменной с помощью функции <code>type()</code>.</p>
        <ul>
            <li><b>Целые числа (<code>int</code>):</b> <code>10</code>, <code>-500</code>.</li>
            <li><b>Числа с плавающей точкой (<code>float</code>):</b> <code>3.14</code>, <code>-0.001</code>.</li>
            <li><b>Строки (<code>str</code>):</b> Текст, заключенный в кавычки (<code>"Python"</code> или <code>''Go''</code>).</li>
            <li><b>Логический (<code>bool</code>):</b> <code>True</code> или <code>False</code>.</li>
        </ul>

        <h2>Арифметические операторы</h2>
        <p>Для работы с числами используются стандартные операторы:</p>
        <ul>
            <li><code>+</code> (сложение)</li>
            <li><code>-</code> (вычитание)</li>
            <li><code>*</code> (умножение)</li>
            <li><code>/</code> (обычное деление, результат <code>float</code>)</li>
            <li><code>//</code> (целочисленное деление, отбрасывает дробную часть)</li>
            <li><code>%</code> (остаток от деления, или оператор "модуль")</li>
            <li><code>**</code> (возведение в степень)</li>
        </ul>
        <pre><code>result_div = 7 / 2    # 3.5
result_int_div = 7 // 2 # 3
result_mod = 7 % 2    # 1 (остаток)</code></pre>
        '
    )
    RETURNING id
),
problem_1_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_1),
        'Привет, мир!',
        'Напишите функцию <code>hello</code>, которая не принимает аргументов и возвращает строку "Hello, World!". Это традиционная первая программа для любого языка.',
        'easy',
        '1 second',
        67108864
    )
    RETURNING id
),
template_1_1 AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_1_1),
        'python',
        'def hello():
    # Ваш код здесь
    return ""
',
        '{ "function_name": "hello", "parameters": [] }'::jsonb
    )
),
tests_1_1 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_1_1),
        '[{"input": {}, "output": "Hello, World!"}]'::jsonb
    )
),

problem_1_2 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_1),
        'Сумма двух чисел',
        'Реализуйте функцию <code>sum_two_numbers</code>, которая принимает два целых числа <code>a</code> и <code>b</code> и возвращает их сумму.',
        'easy',
        '1 second',
        67108864
    )
    RETURNING id
),
template_1_2 AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_1_2),
        'python',
        'def sum_two_numbers(a, b):
    # Ваш код здесь
    return 0
',
        '{ "function_name": "sum_two_numbers", "parameters": [{ "name": "a", "type": "int" }, { "name": "b", "type": "int" }] }'::jsonb
    )
),
tests_1_2 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_1_2),
        '[
            { "input": { "a": 5, "b": 10 }, "output": 15 },
            { "input": { "a": -5, "b": 5 }, "output": 0 },
            { "input": { "a": 0, "b": 0 }, "output": 0 },
            { "input": { "a": 1000, "b": 2000 }, "output": 3000 }
        ]'::jsonb
    )
),

lesson_2 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Условные конструкции: Ветвление логики',
        'Изучение операторов <code>if</code>, <code>elif</code>, <code>else</code> для создания ветвлений, а также логических операторов для сложных условий.',
        2,
        '
        <h1>Урок 2: Условные конструкции: Ветвление логики</h1>
        <p>Условные конструкции позволяют программам принимать решения. В Python это реализуется через блоки <code>if</code>, <code>elif</code> и <code>else</code>. Основой ветвления являются <b>условия</b> — выражения, которые могут быть истинными (<code>True</code>) или ложными (<code>False</code>).</p>

        <h2>Важность отступов (Indentation)</h2>
        <p>В Python отступы (обычно 4 пробела) не просто для красоты, они <b>структурируют код</b>. Блок кода, относящийся к <code>if</code>, <code>elif</code> или <code>else</code>, обязательно должен иметь отступ.</p>

        <h2>Оператор <code>if</code> и <code>else</code></h2>
        <p>Конструкция <code>if</code> выполняет код, если условие истинно. <code>else</code> — выполняет код во всех остальных случаях.</p>
        <pre><code>is_sunny = True
if is_sunny:
    print("Отличный день для прогулки!")
else:
    print("Может быть, стоит остаться дома.")</code></pre>

        <h2>Оператор <code>elif</code> (Else If)</h2>
        <p><code>elif</code> позволяет проверить дополнительное условие, если предыдущее <code>if</code> (или <code>elif</code>) оказалось ложным. Можно использовать любое количество <code>elif</code> между <code>if</code> и <code>else</code>.</p>
        <pre><code>time_of_day = "morning"
if time_of_day == "morning":
    print("Доброе утро!")
elif time_of_day == "afternoon":
    print("Добрый день!")
else:
    print("Добрый вечер!")</code></pre>

        <h2>Операторы сравнения</h2>
        <p>Условия формируются с помощью операторов сравнения:</p>
        <ul>
            <li><code>==</code>: Равно (проверка равенства)</li>
            <li><code>!=</code>: Не равно</li>
            <li><code>></code>: Больше</li>
            <li><code><</code>: Меньше</li>
            <li><code>>=</code>: Больше или равно</li>
            <li><code><=</code>: Меньше или равно</li>
        </ul>

        <h2>Логические операторы</h2>
        <p>Для создания сложных условий используются логические операторы:</p>
        <ul>
            <li><b><code>and</code>:</b> Возвращает <code>True</code>, если оба операнда истинны.</li>
            <li><b><code>or</code>:</b> Возвращает <code>True</code>, если хотя бы один из операндов истинен.</li>
            <li><b><code>not</code>:</b> Инвертирует логическое значение (<code>not True</code> дает <code>False</code>).</li>
        </ul>
        <pre><code>is_old_enough = True
has_permission = False

if is_old_enough and not has_permission:
    print("Доступ ограничен, нужно разрешение.")

if is_old_enough or has_permission:
    print("Хотя бы одно условие выполнено.")</code></pre>
        '
    )
    RETURNING id
),
problem_2_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_2),
        'Чётное или нечётное',
        'Напишите функцию <code>is_even</code>, которая принимает целое число и возвращает <code>True</code>, если оно чётное (делится на 2 без остатка), и <code>False</code> в противном случае. Используйте оператор модуля (<code>%</code>).',
        'easy',
        '1 second',
        67108864
    )
    RETURNING id
),
template_2_1 AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_2_1),
        'python',
        'def is_even(number):
    # Ваш код здесь
    return False
',
        '{ "function_name": "is_even", "parameters": [{ "name": "number", "type": "int" }] }'::jsonb
    )
),
tests_2_1 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_2_1),
        '[
            { "input": { "number": 10 }, "output": true },
            { "input": { "number": 7 }, "output": false },
            { "input": { "number": 0 }, "output": true },
            { "input": { "number": -2 }, "output": true }
        ]'::jsonb
    )
),

problem_2_2 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_2),
        'Високосный год',
        'Напишите функцию <code>is_leap_year</code>, которая определяет, является ли год високосным. Используйте логические операторы и оператор модуля, чтобы проверить правила: год является високосным, если он делится на 4, но не делится на 100, либо если он делится на 400.',
        'medium',
        '1 second',
        67108864
    )
    RETURNING id
),
template_2_2 AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_2_2),
        'python',
        'def is_leap_year(year):
    # Ваш код здесь
    return False
',
        '{ "function_name": "is_leap_year", "parameters": [{ "name": "year", "type": "int" }] }'::jsonb
    )
),
tests_2_2 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_2_2),
        '[
            { "input": { "year": 2024 }, "output": true },
            { "input": { "year": 2023 }, "output": false },
            { "input": { "year": 2000 }, "output": true },
            { "input": { "year": 1900 }, "output": false }
        ]'::jsonb
    )
),

lesson_3 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Циклы: <code>for</code> и <code>while</code> для итераций',
        'Освоение циклов для многократного выполнения кода, а также управление потоком с помощью <code>break</code> и <code>continue</code>.',
        3,
        '
        <h1>Урок 3: Циклы: <code>for</code> и <code>while</code></h1>
        <p>Циклы позволяют автоматизировать повторяющиеся задачи. В Python есть два основных типа циклов, каждый из которых подходит для своей ситуации.</p>

        <h2>Цикл <code>while</code></h2>
        <p>Цикл <code>while</code> ( "пока" ) выполняет блок кода <b>пока</b> его условие истинно. Он идеально подходит, когда вы не знаете заранее, сколько раз потребуется выполнить итерацию.</p>
        <pre><code>i = 1
while i <= 3:
    print(f"Повторение {i}")
    i += 1  # Краткая запись i = i + 1. Обязательно меняйте условие!

# Выведет:
# Повторение 1
# Повторение 2
# Повторение 3</code></pre>
        <p><strong>Внимание:</strong> Всегда следите, чтобы условие цикла <code>while</code> рано или поздно стало <code>False</code>, иначе это приведет к <b>бесконечному циклу</b>, который "заморозит" программу.</p>

        <h2>Цикл <code>for</code></h2>
        <p>Цикл <code>for</code> используется для <b>итерации</b> (перебора) по элементам какой-либо последовательности (строки, списка, диапазона чисел и т.д.).</p>
        
        <h3>Использование <code>range()</code></h3>
        <p>Функция <code>range()</code> генерирует последовательность чисел и часто используется в цикле <code>for</code>:</p>
        <ul>
            <li><code>range(N)</code>: от 0 до N-1.</li>
            <li><code>range(Start, Stop)</code>: от Start до Stop-1.</li>
            <li><code>range(Start, Stop, Step)</code>: от Start до Stop-1 с указанным шагом.</li>
        </ul>
        <pre><code># Итерация 5 раз (i = 0, 1, 2, 3, 4)
for i in range(5):
    print(i)
    
# Итерация по четным числам
for num in range(2, 11, 2):
    print(num) # 2, 4, 6, 8, 10</code></pre>
        
        <h3>Итерация по последовательности</h3>
        <pre><code>my_word = "hello"
for char in my_word:
    print(char.upper()) # H, E, L, L, O
</code></pre>

        <h2>Управление циклом: <code>break</code> и <code>continue</code></h2>
        <ul>
            <li><b><code>break</code>:</b> Немедленно завершает (выходит из) текущего цикла, даже если условие цикла еще истинно.</li>
            <li><b><code>continue</code>:</b> Пропускает оставшуюся часть текущей итерации и переходит к следующей итерации цикла.</li>
        </ul>
        <pre><code># Пример с break
for num in range(10):
    if num == 7:
        break # Цикл остановится на 7
    print(num) # 0, 1, 2, 3, 4, 5, 6</code></pre>
        '
    )
    RETURNING id
),
problem_3_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_3),
        'Сумма чисел до N',
        'Напишите функцию <code>sum_up_to_n</code>, которая вычисляет сумму всех целых чисел от 1 до <code>n</code> включительно. Попробуйте реализовать это с помощью цикла <code>for</code>.',
        'easy',
        '1 second',
        67108864
    )
    RETURNING id
),
template_3_1 AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_3_1),
        'python',
        'def sum_up_to_n(n):
    # Ваш код здесь
    return 0
',
        '{ "function_name": "sum_up_to_n", "parameters": [{ "name": "n", "type": "int" }] }'::jsonb
    )
),
tests_3_1 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_3_1),
        '[
            { "input": { "n": 5 }, "output": 15 },
            { "input": { "n": 1 }, "output": 1 },
            { "input": { "n": 100 }, "output": 5050 }
        ]'::jsonb
    )
),

problem_3_2 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_3),
        'Факториал числа',
        'Напишите функцию <code>factorial</code>, которая вычисляет факториал числа <code>n</code> (произведение всех целых чисел от 1 до <code>n</code>). Факториал 0 равен 1. Используйте цикл <code>while</code> или <code>for</code>.',
        'easy',
        '1 second',
        67108864
    )
    RETURNING id
),
template_3_2 AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_3_2),
        'python',
        'def factorial(n):
    # Ваш код здесь
    return 1
',
        '{ "function_name": "factorial", "parameters": [{ "name": "n", "type": "int" }] }'::jsonb
    )
),
tests_3_2 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_3_2),
        '[
            { "input": { "n": 5 }, "output": 120 },
            { "input": { "n": 0 }, "output": 1 },
            { "input": { "n": 1 }, "output": 1 },
            { "input": { "n": 7 }, "output": 5040 }
        ]'::jsonb
    )
),

lesson_4 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Работа со строками и форматирование',
        'Изучение строк как последовательностей, методы срезов, а также мощные инструменты форматирования, такие как f-строки.',
        4,
        '
        <h1>Урок 4: Работа со строками и форматирование</h1>
        <p>Строка (<code>str</code>) в Python — это упорядоченная последовательность символов, используемая для хранения текстовой информации. Важно помнить, что строки в Python <b>неизменяемы (immutable)</b>. Любая операция, которая кажется изменением, на самом деле создает новую строку.</p>

        <h2>Индексы и срезы (Slicing)</h2>
        <p>Доступ к символам осуществляется по индексу, начиная с <code>0</code>. Отрицательные индексы начинаются с конца (<code>-1</code> — последний символ).</p>
        <p>Срезы позволяют извлечь подстроку с помощью синтаксиса <code>[start:stop:step]</code>.</p>
        <pre><code>my_string = "Программирование"
print(my_string[0])       # П (индекс 0)
print(my_string[-1])      # е (последний символ)
print(my_string[0:4])     # Прог (с 0 по 3)
print(my_string[::2])     # Пргрмиоаи (каждый второй символ)
print(my_string[::-1])    # еинавориммаргорП (перевернутая строка)</code></pre>

        <h2>Важные строковые методы</h2>
        <p>В Python есть богатый набор методов для манипуляций со строками:</p>
        <ul>
            <li><code>.lower()</code> / <code>.upper()</code>: Изменение регистра.</li>
            <li><code>.strip()</code>: Удаление начальных и конечных пробелов.</li>
            <li><code>.replace(old, new)</code>: Замена подстроки.</li>
            <li><code>.split(sep)</code>: Разбиение строки на список по разделителю.</li>
            <li><code>.join(iterable)</code>: Соединение элементов списка в одну строку.</li>
            <li><code>.find(sub)</code>: Поиск индекса первого вхождения подстроки (вернет -1, если не найдено).</li>
            <li><code>.isalpha()</code>, <code>.isdigit()</code>, <code>.isalnum()</code>: Проверки, состоит ли строка только из букв, только из цифр, или из букв и цифр.</li>
        </ul>
        <pre><code>sentence = "Python – лучший. "
words = sentence.split(" ") # ["Python", "–", "лучший.", ""]
new_sentence = " ".join(["Go", "is", "fast"]) # "Go is fast"</code></pre>

        <h2>Форматирование строк (f-строки)</h2>
        <p><b>f-строки</b> (начиная с Python 3.6) — это самый современный и рекомендуемый способ вставки значений переменных в строку. Они позволяют писать выражения прямо внутри фигурных скобок <code>{}</code>.</p>
        <pre><code>item = "ручка"
count = 3
total = 150.50

# Используем f-строку
output = f"У меня есть {count} шт. {item}. Общая стоимость: {total:.2f} руб."
print(output)
# Выведет: У меня есть 3 шт. ручка. Общая стоимость: 150.50 руб.</code></pre>
        '
    )
    RETURNING id
),
problem_4_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_4),
        'Перевернуть строку',
        'Напишите функцию <code>reverse_string</code>, которая принимает строку <code>s</code> и возвращает её в перевёрнутом виде. Попробуйте использовать срез (slicing) <code>[::-1]</code>.',
        'easy',
        '1 second',
        67108864
    )
    RETURNING id
),
template_4_1 AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_4_1),
        'python',
        'def reverse_string(s):
    # Ваш код здесь
    return ""
',
        '{ "function_name": "reverse_string", "parameters": [{ "name": "s", "type": "str" }] }'::jsonb
    )
),
tests_4_1 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_4_1),
        '[
            { "input": { "s": "python" }, "output": "nohtyp" },
            { "input": { "s": "hello" }, "output": "olleh" },
            { "input": { "s": "" }, "output": "" },
            { "input": { "s": "a" }, "output": "a" }
        ]'::jsonb
    )
),
problem_4_2 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_4),
        'Проверка на палиндром',
        'Напишите функцию <code>is_palindrome</code>, которая проверяет, является ли строка <code>s</code> палиндромом. Необходимо игнорировать регистр (привести к нижнему) и удалять все пробелы и знаки препинания (для простоты удаляйте только пробелы).',
        'medium',
        '1 second',
        67108864
    )
    RETURNING id
),
template_4_2 AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_4_2),
        'python',
        'def is_palindrome(s):
    # Ваш код здесь
    return False
',
        '{ "function_name": "is_palindrome", "parameters": [{ "name": "s", "type": "str" }] }'::jsonb
    )
),
tests_4_2 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_4_2),
        '[
            { "input": { "s": "A man a plan a canal Panama" }, "output": true },
            { "input": { "s": "racecar" }, "output": true },
            { "input": { "s": "hello" }, "output": false },
            { "input": { "s": "No lemon, no melon" }, "output": true }
        ]'::jsonb
    )
),

lesson_5 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Списки (Lists): Динамические массивы',
        'Знакомство с одной из самых универсальных структур данных в Python — списками, их изменяемостью и основными методами.',
        5,
        '
        <h1>Урок 5: Списки (Lists): Динамические массивы</h1>
        <p>Список — это упорядоченная, <b>изменяемая (mutable)</b> коллекция элементов. Это означает, что вы можете добавлять, удалять или изменять элементы списка после его создания. Списки создаются с помощью квадратных скобок <code>[]</code>.</p>
        <pre><code>my_list = [10, "apple", 3.14, [1, 2]] # Список может содержать разные типы</code></pre>

        <h2>Доступ и Изменение элементов</h2>
        <p>Элементы доступны по индексу. Благодаря изменяемости, вы можете напрямую менять значения по индексу.</p>
        <pre><code>fruits = ["яблоко", "банан", "апельсин"]
fruits[1] = "груша" # Изменение
print(fruits)       # ["яблоко", "груша", "апельсин"]</code></pre>

        <h2>Основные методы списков</h2>
        <p>Методы изменяют список <b>на месте</b> (in-place) и не возвращают новый список (кроме <code>.pop()</code>, который возвращает удаленный элемент).</p>
        <ul>
            <li><code>.append(element)</code>: Добавляет элемент в конец.</li>
            <li><code>.insert(index, element)</code>: Вставляет элемент по указанному индексу.</li>
            <li><code>.pop([index])</code>: Удаляет и возвращает элемент по индексу (по умолчанию — последний).</li>
            <li><code>.remove(value)</code>: Удаляет первое вхождение указанного значения.</li>
            <li><code>.sort()</code>: Сортирует список.</li>
            <li><code>.reverse()</code>: Разворачивает список.</li>
        </ul>
        <pre><code>data = [3, 1, 2]
data.append(4)      # [3, 1, 2, 4]
data.sort()         # [1, 2, 3, 4]
data.remove(2)      # [1, 3, 4]
last = data.pop()   # last = 4, data = [1, 3]</code></pre>

        <h2>Итерация по спискам</h2>
        <p>Используйте цикл <code>for</code> для перебора элементов. Если вам нужен и индекс, и значение, используйте функцию <code>enumerate()</code>.</p>
        <pre><code>names = ["Alex", "Bob", "Chris"]
for index, name in enumerate(names):
    print(f"Имя {index}: {name}")
    
# Выведет:
# Имя 0: Alex
# Имя 1: Bob
# Имя 2: Chris</code></pre>
        '
    )
    RETURNING id
),
problem_5_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_5),
        'Найти максимальный элемент',
        'Напишите функцию <code>find_max</code>, которая принимает список чисел <code>numbers</code> и возвращает максимальное значение, не используя встроенную функцию <code>max()</code>. Если список пуст, верните <code>None</code>.',
        'easy',
        '1 second',
        67108864
    )
    RETURNING id
),
template_5_1 AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_5_1),
        'python',
        'def find_max(numbers):
    # Ваш код здесь
    return None
',
        '{ "function_name": "find_max", "parameters": [{ "name": "numbers", "type": "list" }] }'::jsonb
    )
),
tests_5_1 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_5_1),
        '[
            { "input": { "numbers": [1, 5, 2, 9, 3] }, "output": 9 },
            { "input": { "numbers": [-1, -5, -2] }, "output": -1 },
            { "input": { "numbers": [10] }, "output": 10 },
            { "input": { "numbers": [] }, "output": null }
        ]'::jsonb
    )
),
problem_5_2 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_5),
        'Удаление дубликатов',
        'Напишите функцию <code>remove_duplicates</code>, которая принимает список <code>items</code> и возвращает новый список, содержащий только уникальные элементы из исходного, <b>сохраняя их первоначальный порядок</b>.',
        'medium',
        '1 second',
        67108864
    )
    RETURNING id
),
template_5_2 AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_5_2),
        'python',
        'def remove_duplicates(items):
    # Ваш код здесь
    return []
',
        '{ "function_name": "remove_duplicates", "parameters": [{ "name": "items", "type": "list" }] }'::jsonb
    )
),
tests_5_2 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_5_2),
        '[
            { "input": { "items": [1, 2, 2, 3, 1, 4] }, "output": [1, 2, 3, 4] },
            { "input": { "items": ["a", "b", "a", "c"] }, "output": ["a", "b", "c"] },
            { "input": { "items": [10, 10, 10, 10] }, "output": [10] },
            { "input": { "items": [] }, "output": [] }
        ]'::jsonb
    )
),

lesson_6 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Словари (Dictionaries): Коллекции "ключ-значение"',
        'Изучение словарей для хранения данных в формате "ключ-значение", их эффективного использования и методов итерации.',
        6,
        '
        <h1>Урок 6: Словари (Dictionaries): Коллекции "ключ-значение"</h1>
        <p>Словарь (<code>dict</code>) — это неупорядоченная, <b>изменяемая</b> коллекция, где элементы хранятся в виде пар <b>"ключ: значение"</b>. Словари оптимизированы для очень быстрого поиска, получения и удаления данных по ключу. Ключи должны быть <b>уникальными</b> и <b>хешируемыми</b> (неизменяемыми, например, строки или числа).</p>
        <p>Словари создаются с помощью фигурных скобок <code>{}</code>.</p>
        <pre><code>student = {
    "id": 101,
    "name": "Олег",
    "scores": [90, 85, 95]
}</code></pre>

        <h2>Доступ, добавление и изменение</h2>
        <p>Доступ и присваивание происходят по ключу в квадратных скобках <code>[]</code>. Если ключа нет, присваивание добавляет новую пару, а попытка доступа с помощью <code>[]</code> вызовет ошибку <code>KeyError</code>.</p>
        <pre><code># Доступ
print(student["name"])  # Выведет "Олег"

# Добавление/Изменение
student["name"] = "Игорь"
student["city"] = "Минск"

# Безопасный доступ с .get()
age = student.get("age", "Неизвестно") # "Неизвестно"</code></pre>

        <h2>Удаление элементов</h2>
        <ul>
            <li><code>del my_dict[key]</code>: Удаляет пару по ключу.</li>
            <li><code>.pop(key, [default])</code>: Удаляет и возвращает значение по ключу.</li>
        </ul>

        <h2>Итерация по словарю</h2>
        <p>Вы можете перебирать ключи, значения или пары ключ-значение:</p>
        <pre><code># 1. Перебор ключей (по умолчанию)
for key in student:
    print(key)

# 2. Перебор значений
for value in student.values():
    print(value)

# 3. Перебор пар (самый частый способ)
for key, value in student.items():
    print(f"{key}: {value}")</code></pre>
        '
    )
    RETURNING id
),
problem_6_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_6),
        'Подсчёт частоты слов',
        'Напишите функцию <code>word_frequency</code>, которая принимает строку (текст) и возвращает словарь, где ключи — это слова, а значения — количество их повторений. Перед подсчетом необходимо: 1) удалить знаки препинания (для простоты удаляйте только <code>.</code>, <code>,</code>, <code>!</code>, <code>?</code>); 2) привести все слова к нижнему регистру.',
        'medium',
        '2 seconds',
        134217728
    )
    RETURNING id
),
template_6_1 AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_6_1),
        'python',
        'def word_frequency(text):
    # Ваш код здесь
    return {}
',
        '{ "function_name": "word_frequency", "parameters": [{ "name": "text", "type": "str" }] }'::jsonb
    )
)
INSERT INTO tests (problem_id, tests)
VALUES (
    (SELECT id FROM problem_6_1),
    '[
        { "input": { "text": "hello world hello" }, "output": { "hello": 2, "world": 1 } },
        { "input": { "text": "Python is awesome and Python is fun" }, "output": { "python": 2, "is": 2, "awesome": 1, "and": 1, "fun": 1 } },
        { "input": { "text": "Go, go, go! It is time to go." }, "output": { "go": 4, "it": 1, "is": 1, "time": 1, "to": 1 } },
        { "input": { "text": "A a B b c c A" }, "output": { "a": 2, "b": 2, "c": 2 } }
    ]'::jsonb
);

COMMIT;
-- +goose StatementEnd
