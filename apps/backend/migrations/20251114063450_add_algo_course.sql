-- +goose Up
-- +goose StatementBegin
BEGIN;

WITH course_insert AS (
    INSERT INTO courses (name, description, tags, status)
    VALUES (
        'Основы алгоритмов и структур данных',
        'Курс, посвященный фундаментальным алгоритмическим концепциям и простейшим структурам данных. Цель — научиться писать не просто правильный, но и эффективный код.',
        '{"algorithms", "data-structures", "complexity", "general"}',
        'published'
    )
    RETURNING id
),

lesson_1 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Введение в алгоритмы и Big O-нотация',
        'Понимание того, как оценивать эффективность кода с помощью Big O-нотации (временная и пространственная сложность).',
        1,
        $$
        <h1>Урок 1: Введение в алгоритмы и Big O-нотация</h1>
        <p><b>Алгоритм</b> — это пошаговая инструкция для решения определенной задачи. В программировании важно не только, чтобы код работал, но и то, насколько эффективно (быстро и с малым потреблением памяти) он это делает.</p>

        <h2>Big O-нотация (О-большое)</h2>
        <p>Big O-нотация — это математический способ описания <b>верхней границы</b> временной (time) или пространственной (space) сложности алгоритма. Она показывает, как время выполнения или объем памяти растет по мере увеличения размера входных данных (N).</p>

        <h3>Типичные сложности:</h3>
        <ul>
            <li><b>O(1) - Константная:</b> Операция выполняется за одно и то же время, независимо от N (например, доступ к элементу массива по индексу).</li>
            <li><b>O(log N) - Логарифмическая:</b> Время выполнения растет очень медленно. Типично для алгоритмов "разделяй и властвуй" (например, бинарный поиск).</li>
            <li><b>O(N) - Линейная:</b> Время выполнения прямо пропорционально N (например, перебор всех элементов списка).</li>
            <li><b>O(N log N) - Линейно-логарифмическая:</b> Типично для эффективных алгоритмов сортировки (например, QuickSort, MergeSort).</li>
            <li><b>O(N²) - Квадратичная:</b> Время выполнения растет квадратично. Типично для циклов, вложенных друг в друга (например, простая сортировка пузырьком).</li>
            <li><b>O(2ⁿ) - Экспоненциальная:</b> Очень медленно. Подходит только для очень маленьких N (например, рекурсивное вычисление чисел Фибоначчи без оптимизации).</li>
        </ul>
        <pre><code>// O(1) - Константная сложность (Go)
func getFirst(arr []int) int {
    return arr[0]
}

# O(N) - Линейная сложность (Python)
def sum_array(arr):
    sum = 0
    for val in arr:
        sum += val
    return sum</code></pre>
        $$
    )
    RETURNING id
),
problem_1_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_1),
        'Поиск среднего',
        'Напишите функцию <code>Average</code> (Go) / <code>average</code> (Python), которая принимает массив/список целых чисел и возвращает их среднее арифметическое (с плавающей точкой). Оцените сложность по времени.',
        'easy',
        '1 second',
        67108864
    )
    RETURNING id
),
template_1_1_go AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_1_1),
        'go',
        $$
// Average вычисляет среднее арифметическое элементов слайса.
func Average(nums []int) float64 {
    // Ваша реализация. Помните о преобразовании типов при делении.
    return 0.0
}
$$,
        '{"function_name": "Average", "parameters": [{"name": "nums", "type": "[]int"}], "return_type": "float64"}'::jsonb
    )
),
template_1_1_python AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_1_1),
        'python',
        $$
def average(nums):
    # Ваша реализация. Помните о делении с плавающей точкой.
    return 0.0
$$,
        '{"function_name": "average", "parameters": [{"name": "nums", "type": "list"}], "return_type": "float"}'::jsonb
    )
),
tests_1_1 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_1_1),
        '[
            { "input": { "nums": [1, 2, 3, 4, 5] }, "output": 3.0 },
            { "input": { "nums": [10, 20] }, "output": 15.0 },
            { "input": { "nums": [] }, "output": 0.0 }
        ]'::jsonb
    )
),

problem_1_2 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_1),
        'Константный доступ',
        'Напишите функцию <code>GetFirst</code>, которая принимает непустой массив/список целых чисел и возвращает его первый элемент. Это должно быть выполнено за O(1).',
        'easy',
        '1 second',
        67108864
    )
    RETURNING id
),
template_1_2_go AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_1_2),
        'go',
        $$
// GetFirst возвращает первый элемент слайса.
func GetFirst(nums []int) int {
    // Вставьте ваш код
    return 0
}
$$,
        '{"function_name": "GetFirst", "parameters": [{"name": "nums", "type": "[]int"}], "return_type": "int"}'::jsonb
    )
),
template_1_2_python AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_1_2),
        'python',
        $$
def get_first(nums):
    # Вставьте ваш код
    return 0
$$,
        '{"function_name": "get_first", "parameters": [{"name": "nums", "type": "list"}], "return_type": "int"}'::jsonb
    )
),
tests_1_2 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_1_2),
        '[
            { "input": { "nums": [10, 20, 30] }, "output": 10 },
            { "input": { "nums": [5] }, "output": 5 }
        ]'::jsonb
    )
),

lesson_2 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Поиск: Линейный и Бинарный',
        'Изучение двух основных алгоритмов поиска и их сложности: O(N) и O(log N).',
        2,
        $$
        <h1>Урок 2: Поиск: Линейный и Бинарный</h1>
        <p>Поиск элемента в коллекции является одной из самых частых задач в программировании. Эффективность поиска зависит от того, как организованы данные.</p>

        <h2>Линейный поиск (Linear Search)</h2>
        <p>Простейший метод: последовательный перебор каждого элемента, пока не будет найден искомый. Работает на любых коллекциях.</p>
        <ul>
            <li><b>Сложность:</b> O(N) — в худшем случае нужно проверить все N элементов.</li>
        </ul>
        <pre><code>// Go (Линейный поиск)
func LinearSearch(arr []int, target int) bool {
    for _, val := range arr {
        if val == target {
            return true
        }
    }
    return false
}</code></pre>

        <h2>Бинарный поиск (Binary Search)</h2>
        <p>Бинарный поиск — это гораздо более быстрый алгоритм, но у него есть строгое требование: коллекция должна быть <b>отсортирована</b>.</p>
        <p>Алгоритм работает так:
            <ol>
                <li>Находим элемент в середине.</li>
                <li>Если середина — искомый элемент, возвращаем его.</li>
                <li>Если искомый элемент меньше середины, ищем только в левой половине.</li>
                <li>Если больше — ищем только в правой половине.</li>
            </ol>
        <p>Таким образом, мы на каждом шаге отбрасываем половину оставшихся данных.</p>
        <ul>
            <li><b>Сложность:</b> O(log N) — очень быстро для больших коллекций.</li>
        </ul>
        $$
    )
    RETURNING id
),
problem_2_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_2),
        'Линейный поиск',
        'Реализуйте функцию <code>LinearSearch</code> (Go) / <code>linear_search</code> (Python), которая принимает отсортированный или неотсортированный массив/список чисел <code>nums</code> и целевое число <code>target</code>. Функция должна вернуть <code>true</code>, если <code>target</code> найдено, и <code>false</code> в противном случае.',
        'easy',
        '1 second',
        67108864
    )
    RETURNING id
),
template_2_1_go AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_2_1),
        'go',
        $$
// LinearSearch ищет элемент в слайсе.
func LinearSearch(nums []int, target int) bool {
    // Вставьте ваш код
    return false
}
$$,
        '{"function_name": "LinearSearch", "parameters": [{"name": "nums", "type": "[]int"}, {"name": "target", "type": "int"}], "return_type": "bool"}'::jsonb
    )
),
template_2_1_python AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_2_1),
        'python',
        $$
def linear_search(nums, target):
    # Вставьте ваш код
    return False
$$,
        '{"function_name": "linear_search", "parameters": [{"name": "nums", "type": "list"}, {"name": "target", "type": "int"}], "return_type": "bool"}'::jsonb
    )
),
tests_2_1 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_2_1),
        '[
            { "input": { "nums": [5, 2, 8, 1, 9], "target": 8 }, "output": true },
            { "input": { "nums": [5, 2, 8, 1, 9], "target": 7 }, "output": false },
            { "input": { "nums": [], "target": 1 }, "output": false }
        ]'::jsonb
    )
),

problem_2_2 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_2),
        'Бинарный поиск',
        'Реализуйте функцию <code>BinarySearch</code> (Go) / <code>binary_search</code> (Python), которая принимает <b>отсортированный</b> массив/список чисел <code>nums</code> и целевое число <code>target</code>. Функция должна вернуть индекс <code>target</code>, если оно найдено, и <code>-1</code> в противном случае. Сложность должна быть O(log N).',
        'medium',
        '1 second',
        67108864
    )
    RETURNING id
),
template_2_2_go AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_2_2),
        'go',
        $$
// BinarySearch ищет элемент в отсортированном слайсе, возвращая его индекс или -1.
func BinarySearch(nums []int, target int) int {
    // Вставьте ваш код
    return -1
}
$$,
        '{"function_name": "BinarySearch", "parameters": [{"name": "nums", "type": "[]int"}, {"name": "target", "type": "int"}], "return_type": "int"}'::jsonb
    )
),
template_2_2_python AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_2_2),
        'python',
        $$
def binary_search(nums, target):
    # Вставьте ваш код
    return -1
$$,
        '{"function_name": "binary_search", "parameters": [{"name": "nums", "type": "list"}, {"name": "target", "type": "int"}], "return_type": "int"}'::jsonb
    )
),
tests_2_2 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_2_2),
        '[
            { "input": { "nums": [2, 5, 8, 12, 16], "target": 12 }, "output": 3 },
            { "input": { "nums": [2, 5, 8, 12, 16], "target": 7 }, "output": -1 },
            { "input": { "nums": [5], "target": 5 }, "output": 0 },
            { "input": { "nums": [1, 3, 5, 7, 9], "target": 1 }, "output": 0 }
        ]'::jsonb
    )
),

lesson_3 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Сортировка: Пузырьком',
        'Изучение одного из простейших, но наименее эффективных алгоритмов сортировки — Сортировка пузырьком (Bubble Sort).',
        3,
        $$
        <h1>Урок 3: Сортировка: Пузырьком (Bubble Sort)</h1>
        <p>Сортировка — это процесс упорядочивания элементов. Сортировка пузырьком — это самый простой алгоритм, который показывает, как работают основные принципы обмена элементами.</p>

        <h2>Принцип работы Bubble Sort</h2>
        <p>Алгоритм многократно проходит по списку (массиву/слайсу), сравнивая соседние элементы и меняя их местами, если они стоят в неправильном порядке. С каждым проходом самое большое "всплывает" в конец, как пузырек.</p>
        <p>Для списка из N элементов требуется:
            <ol>
                <li>Внешний цикл, который повторяется N-1 раз.</li>
                <li>Внутренний цикл, который выполняет сравнения и обмены.</li>
            </ol>
        </p>

        <h2>Сложность</h2>
        <ul>
            <li><b>Временная сложность:</b> O(N²) — два вложенных цикла, делающих почти N*N сравнений. Это делает его неэффективным для больших объемов данных.</li>
            <li><b>Пространственная сложность:</b> O(1) — сортировка происходит "на месте", не требуется дополнительной памяти.</li>
        </ul>
        <pre><code>// Псевдокод
function bubbleSort(array):
    n = length(array)
    for i from 0 to n-2:
        for j from 0 to n-2-i:
            if array[j] > array[j+1]:
                swap(array[j], array[j+1])</code></pre>
        $$
    )
    RETURNING id
),
problem_3_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_3),
        'Сортировка пузырьком',
        'Реализуйте функцию <code>BubbleSort</code> (Go) / <code>bubble_sort</code> (Python), которая принимает массив/список целых чисел <code>nums</code> и сортирует его <b>по возрастанию</b> "на месте" с помощью алгоритма Сортировки пузырьком.',
        'medium',
        '2 seconds',
        67108864
    )
    RETURNING id
),
template_3_1_go AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_3_1),
        'go',
        $$
// BubbleSort сортирует слайс целых чисел по возрастанию.
func BubbleSort(nums []int) {
    // Вставьте ваш код. В Go используйте nums[j], nums[j+1] = nums[j+1], nums[j] для обмена.
}
$$,
        '{"function_name": "BubbleSort", "parameters": [{"name": "nums", "type": "[]int"}], "return_type": "void"}'::jsonb
    )
),
template_3_1_python AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_3_1),
        'python',
        $$
def bubble_sort(nums):
    # Вставьте ваш код. Функция должна изменять список "на месте".
    pass
$$,
        '{"function_name": "bubble_sort", "parameters": [{"name": "nums", "type": "list"}], "return_type": "void"}'::jsonb
    )
),
tests_3_1 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_3_1),
        '[
            { "input": { "nums": [5, 1, 4, 2, 8] }, "output": [1, 2, 4, 5, 8] },
            { "input": { "nums": [1, 2, 3] }, "output": [1, 2, 3] },
            { "input": { "nums": [3, 2, 1] }, "output": [1, 2, 3] },
            { "input": { "nums": [] }, "output": [] }
        ]'::jsonb
    )
),

lesson_4 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Рекурсия: Самовызывающиеся функции',
        'Изучение принципов рекурсии, базового случая и ее применения в алгоритмах.',
        4,
        $$
        <h1>Урок 4: Рекурсия</h1>
        <p><b>Рекурсия</b> — это процесс, при котором функция вызывает сама себя. Рекурсивное решение часто более элегантно, чем итеративное, для задач, которые естественно распадаются на меньшие, похожие подзадачи (например, обход дерева, факториал).</p>

        <h2>Три правила рекурсии</h2>
        <ol>
            <li><b>Базовый случай (Base Case):</b> Должен быть условие, при котором функция прекращает вызывать себя. Это предотвращает бесконечную рекурсию.</li>
            <li><b>Рекурсивный вызов:</b> Функция должна вызывать сама себя, но для меньшей (более простой) подзадачи.</li>
            <li><b>Прогресс:</b> Каждый рекурсивный вызов должен приближать функцию к базовому случаю.</li>
        </ol>

        <h2>Пример: Факториал</h2>
        <p>Факториал числа N (N!) — это произведение всех целых чисел от 1 до N. Рекурсивное определение: <code>N! = N * (N-1)!</code>, а базовый случай: <code>0! = 1</code>.</p>
        <pre><code>// Python: Рекурсивный факториал
def factorial(n):
    if n == 0:  # Базовый случай
        return 1
    return n * factorial(n - 1) # Рекурсивный вызов</code></pre>
        <p>Рекурсия часто имеет пространственную сложность O(N) из-за необходимости хранить все вызовы функций в стеке.</p>
        $$
    )
    RETURNING id
),
problem_4_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_4),
        'Рекурсивный факториал',
        'Реализуйте функцию <code>Factorial</code> (Go) / <code>factorial</code> (Python), которая вычисляет факториал числа <code>n</code>, используя <b>рекурсию</b>.',
        'easy',
        '1 second',
        67108864
    )
    RETURNING id
),
template_4_1_go AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_4_1),
        'go',
        $$
// Factorial вычисляет факториал числа n рекурсивно.
func Factorial(n int) int {
    // Вставьте ваш код
    return 0
}
$$,
        '{"function_name": "Factorial", "parameters": [{"name": "n", "type": "int"}], "return_type": "int"}'::jsonb
    )
),
template_4_1_python AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_4_1),
        'python',
        $$
def factorial(n):
    # Вставьте ваш код
    return 0
$$,
        '{"function_name": "factorial", "parameters": [{"name": "n", "type": "int"}], "return_type": "int"}'::jsonb
    )
),
tests_4_1 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_4_1),
        '[
            { "input": { "n": 5 }, "output": 120 },
            { "input": { "n": 0 }, "output": 1 },
            { "input": { "n": 1 }, "output": 1 },
            { "input": { "n": 7 }, "output": 5040 }
        ]'::jsonb
    )
),

problem_4_2 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_4),
        'Числа Фибоначчи',
        'Напишите функцию <code>Fibonacci</code> (Go) / <code>fibonacci</code> (Python), которая возвращает N-е число Фибоначчи, используя <b>рекурсию</b>. (F(0)=0, F(1)=1, F(N) = F(N-1) + F(N-2)).',
        'medium',
        '2 seconds',
        67108864
    )
    RETURNING id
),
template_4_2_go AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_4_2),
        'go',
        $$
// Fibonacci возвращает n-е число Фибоначчи рекурсивно.
func Fibonacci(n int) int {
    // Вставьте ваш код
    return 0
}
$$,
        '{"function_name": "Fibonacci", "parameters": [{"name": "n", "type": "int"}], "return_type": "int"}'::jsonb
    )
),
template_4_2_python AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_4_2),
        'python',
        $$
def fibonacci(n):
    # Вставьте ваш код
    return 0
$$,
        '{"function_name": "fibonacci", "parameters": [{"name": "n", "type": "int"}], "return_type": "int"}'::jsonb
    )
),
tests_4_2 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_4_2),
        '[
            { "input": { "n": 0 }, "output": 0 },
            { "input": { "n": 1 }, "output": 1 },
            { "input": { "n": 2 }, "output": 1 },
            { "input": { "n": 6 }, "output": 8 }
        ]'::jsonb
    )
),

lesson_5 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Хеширование и Словари/Карты',
        'Использование встроенных хеш-таблиц (<code>map</code> в Go, <code>dict</code> в Python) для решения задач с O(1) доступом.',
        5,
        $$
        <h1>Урок 5: Хеширование и Словари/Карты</h1>
        <p>Наиболее эффективным способом хранения и поиска данных является использование <b>хеш-таблиц</b> (<code>map</code> в Go, <code>dict</code> в Python). В идеале, доступ, вставка и удаление элемента в хеш-таблице имеют сложность <b>O(1)</b>.</p>

        <h2>Применение: Подсчет частоты</h2>
        <p>Хеш-таблицы идеально подходят для подсчета частоты элементов. Вместо того, чтобы проходить по списку несколько раз (что дало бы O(N²) или O(N log N)), мы проходим по нему только один раз (O(N)), используя хеш-таблицу для хранения счетчиков.</p>
        <pre><code>// Python: Подсчет частоты
counts = {}
for item in my_list:
    counts[item] = counts.get(item, 0) + 1</code></pre>
        
        <h2>Задача двух сумм</h2>
        <p>Классическая задача: найти два числа в списке, сумма которых равна целевому значению.
        <ul>
            <li><b>Наивный (O(N²)):</b> Использовать два вложенных цикла, чтобы перебрать все пары.</li>
            <li><b>Эффективный (O(N)):</b> Использовать хеш-таблицу. Для каждого числа <code>x</code> ищем, есть ли в таблице <code>target - x</code>.</li>
        </ul>
        $$
    )
    RETURNING id
),
problem_5_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_5),
        'Две суммы (Two Sum)',
        'Напишите функцию <code>TwoSum</code> (Go) / <code>two_sum</code> (Python), которая принимает массив/список целых чисел <code>nums</code> и целевое число <code>target</code>. Функция должна вернуть индексы двух чисел, сумма которых равна <code>target</code>. Предполагается, что существует ровно одно решение. Используйте хеш-таблицу для достижения <b>O(N)</b> сложности.',
        'medium',
        '1 second',
        67108864
    )
    RETURNING id
),
template_5_1_go AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_5_1),
        'go',
        $$
// TwoSum находит индексы двух чисел, сумма которых равна target.
func TwoSum(nums []int, target int) []int {
    // Вставьте ваш код: map для хранения [значение]индекс
    return []int{}
}
$$,
        '{"function_name": "TwoSum", "parameters": [{"name": "nums", "type": "[]int"}, {"name": "target", "type": "int"}], "return_type": "[]int"}'::jsonb
    )
),
template_5_1_python AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_5_1),
        'python',
        $$
def two_sum(nums, target):
    # Вставьте ваш код: словарь для хранения {значение: индекс}
    return []
$$,
        '{"function_name": "two_sum", "parameters": [{"name": "nums", "type": "list"}, {"name": "target", "type": "int"}], "return_type": "list"}'::jsonb
    )
),
tests_5_1 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_5_1),
        '[
            { "input": { "nums": [2, 7, 11, 15], "target": 9 }, "output": [0, 1] },
            { "input": { "nums": [3, 2, 4], "target": 6 }, "output": [1, 2] },
            { "input": { "nums": [3, 3], "target": 6 }, "output": [0, 1] }
        ]'::jsonb
    )
),

lesson_6 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Алгоритмы для строк: Скользящее окно и Двух указателей',
        'Изучение двух популярных методов для оптимизации строковых алгоритмов: "Скользящее окно" и "Два указателя".',
        6,
        $$
        <h1>Урок 6: Алгоритмы для строк</h1>
        <p>Работа со строками часто требует более умных подходов, чем просто вложенные циклы. Два эффективных шаблона для оптимизации — это <b>Два указателя</b> и <b>Скользящее окно</b>.</p>

        <h2>Шаблон "Два указателя" (Two Pointers)</h2>
        <p>Этот шаблон использует два указателя (индекса), которые движутся по массиву или строке в одном направлении (например, оба от начала) или в противоположных (один от начала, другой от конца). Это часто позволяет решать задачи с O(N²) до O(N).</p>
        <p><b>Типичное применение:</b> Проверка палиндромов, поиск пар в отсортированном массиве.</p>
        <pre><code>// Псевдокод для проверки палиндрома
left = 0
right = length(s) - 1
while left < right:
    if s[left] != s[right]:
        return false
    left += 1
    right -= 1
return true</code></pre>

        <h2>Шаблон "Скользящее окно" (Sliding Window)</h2>
        <p>Скользящее окно — это подмассив или подстрока в основном массиве/строке, которая перемещается от начала до конца. Это очень полезно для задач, где нужно найти максимальную/минимальную сумму, длину или количество в <b>подпоследовательности</b> заданной длины или с определенным условием.</p>
        <p>Окно "скользит", расширяясь с одной стороны (<code>right</code>) и сужаясь с другой (<code>left</code>), чтобы соответствовать условию. Это позволяет избежать O(N²) перепроверок всех подстрок, сводя сложность к O(N).</p>
        $$
    )
    RETURNING id
),
problem_6_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_6),
        'Обратный порядок слов',
        'Напишите функцию <code>ReverseWords</code> (Go) / <code>reverse_words</code> (Python), которая принимает строку, состоящую из слов, разделенных пробелами. Функция должна вернуть строку, в которой слова расположены в обратном порядке (например, "the sky is blue" -> "blue is sky the").',
        'medium',
        '1 second',
        67108864
    )
    RETURNING id
),
template_6_1_go AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_6_1),
        'go',
        $$
// ReverseWords возвращает строку с обратным порядком слов.
func ReverseWords(s string) string {
    // Вставьте ваш код. Полезно использовать strings.Fields()
    return ""
}
$$,
        '{"function_name": "ReverseWords", "parameters": [{"name": "s", "type": "string"}], "return_type": "string"}'::jsonb
    )
),
template_6_1_python AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_6_1),
        'python',
        $$
def reverse_words(s):
    # Вставьте ваш код. Полезно использовать s.split() и " ".join()
    return ""
$$,
        '{"function_name": "reverse_words", "parameters": [{"name": "s", "type": "str"}], "return_type": "str"}'::jsonb
    )
),
tests_6_1 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_6_1),
        '[
            { "input": { "s": "the sky is blue" }, "output": "blue is sky the" },
            { "input": { "s": "  hello world  " }, "output": "world hello" },
            { "input": { "s": "a single word" }, "output": "word single a" }
        ]'::jsonb
    )
),

problem_6_2 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_6),
        'Максимальный подмассив с суммой K',
        'Напишите функцию <code>MaxSubarrayLen</code> (Go) / <code>max_subarray_len</code> (Python), которая находит длину <b>самого длинного подмассива</b>, сумма элементов которого равна заданному числу <code>k</code>. Используйте метод **Скользящее окно** или хеш-таблицу.',
        'hard',
        '1 second',
        67108864
    )
    RETURNING id
),
template_6_2_go AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_6_2),
        'go',
        $$
// MaxSubarrayLen находит длину самого длинного подмассива с суммой k.
func MaxSubarrayLen(nums []int, k int) int {
    // Вставьте ваш код
    return 0
}
$$,
        '{"function_name": "MaxSubarrayLen", "parameters": [{"name": "nums", "type": "[]int"}, {"name": "k", "type": "int"}], "return_type": "int"}'::jsonb
    )
),
template_6_2_python AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_6_2),
        'python',
        $$
def max_subarray_len(nums, k):
    # Вставьте ваш код
    return 0
$$,
        '{"function_name": "max_subarray_len", "parameters": [{"name": "nums", "type": "list"}, {"name": "k", "type": "int"}], "return_type": "int"}'::jsonb
    )
)
INSERT INTO tests (problem_id, tests)
VALUES (
    (SELECT id FROM problem_6_2),
    '[
        { "input": { "nums": [1, -1, 5, -2, 3], "k": 3 }, "output": 4 }, 
        { "input": { "nums": [-2, -1, 2, 1], "k": 1 }, "output": 2 },    
        { "input": { "nums": [1, 2, 3], "k": 7 }, "output": 0 }
    ]'::jsonb
);

COMMIT;
-- +goose StatementEnd
