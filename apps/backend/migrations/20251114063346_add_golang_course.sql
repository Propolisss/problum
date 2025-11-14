-- +goose Up
-- +goose StatementBegin
BEGIN;

WITH course_insert AS (
    INSERT INTO courses (name, description, tags, status)
    VALUES (
        'Основы Go: Структура, Типы и Функции',
        'Комплексный курс для изучения основ языка программирования Go. Мы рассмотрим его уникальный синтаксис, строгую типизацию, работу с функциями, массивами, слайсами и картами.',
        '{"go", "golang", "beginner", "concurrency"}',
        'published'
    )
    RETURNING id
),

lesson_1 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Структура программы, переменные и типы',
        'Знакомство с базовой структурой программы Go, строгой типизацией, объявлением переменных и основными типами данных.',
        1,
        $$
        <h1>Урок 1: Структура Go, переменные и типы</h1>
        <p>Go — это скомпилированный, статически типизированный язык. Каждая программа на Go состоит из пакетов, а выполнение начинается с функции <code>main</code> в пакете <code>main</code>.</p>

        <h2>Базовая структура программы</h2>
        <p>Любая исполняемая программа на Go начинается с объявления пакета <code>main</code> и функции <code>main()</code>.</p>
        <pre><code>package main

import "fmt" // Импорт пакета для форматированного ввода/вывода

func main() {
    // Весь исполняемый код начинается отсюда
    fmt.Println("Hello, Go!")
}</code></pre>
        <p><code>import</code> используется для подключения других пакетов, таких как <code>fmt</code> (от "format"), который содержит функцию <code>Println</code> для вывода в консоль.</p>

        <h2>Комментарии и Точка с запятой</h2>
        <p>Комментарии, как и в Python, начинаются с <code>//</code>. В Go, в отличие от C/C++, вам не нужно ставить точку с запятой (<code>;</code>) в конце большинства операторов. Компилятор Go вставляет их автоматически.</p>

        <h2>Объявление переменных и типы</h2>
        <p>Go — <b>статически типизированный</b> язык, что означает, что тип переменной должен быть известен во время компиляции. Существует три основных способа объявления:</p>
        <ol>
            <li><b>С явным указанием типа:</b>
                <pre><code>var age int = 30
var name string = "Алекс"</code></pre>
            </li>
            <li><b>С выводом типа (Type Inference):</b>
                <pre><code>var count = 50   // Go сам определит, что это int
var price = 99.99 // Go сам определит, что это float64</code></pre>
            </li>
            <li><b>Короткое объявление (самый популярный способ):</b> Используется оператор <code>:=</code>. Работает только внутри функций.
                <pre><code>isTrue := true
message := "Go rocks!"</code></pre>
            </li>
        </ol>

        <h2>Основные типы данных</h2>
        <ul>
            <li><b>Числа:</b>
                <ul>
                    <li><code>int</code>, <code>int8</code>, <code>int16</code>, <code>int32</code>, <code>int64</code> (целые числа)</li>
                    <li><code>float32</code>, <code>float64</code> (числа с плавающей точкой)</li>
                </ul>
            </li>
            <li><b>Строки:</b> <code>string</code>.</li>
            <li><b>Логический:</b> <code>bool</code> (<code>true</code> или <code>false</code>).</li>
        </ul>
        <p>Переменная, объявленная, но не инициализированная, получает <b>нулевое значение</b> (zero value) своего типа: <code>0</code> для чисел, <code>""</code> для строк, <code>false</code> для булевых.</p>
        $$
    )
    RETURNING id
),
problem_1_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_1),
        'Сумма двух чисел',
        'Напишите функцию <code>Sum</code>, которая принимает два целых числа <code>a</code> и <code>b</code> типа <code>int</code> и возвращает их сумму.',
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
// Sum принимает два целых числа и должна вернуть их сумму.
func Sum(a int, b int) int {
    // Вставьте ваш код здесь
}
$$,
        '{"function_name": "Sum", "parameters": [{"name": "a", "type": "int"}, {"name": "b", "type": "int"}]}'::jsonb
    )
),
tests_1_1 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_1_1),
        '[
          { "input": { "a": 1, "b": 2 }, "output": 3 },
          { "input": { "a": -5, "b": 5 }, "output": 0 },
          { "input": { "a": 0, "b": 0 }, "output": 0 },
          { "input": { "a": 1000, "b": 2000 }, "output": 3000 }
        ]'::jsonb
    )
),

problem_1_2 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_1),
        'Площадь прямоугольника',
        'Напишите функцию <code>Area</code>, которая принимает длину <code>length</code> и ширину <code>width</code> (оба <code>float64</code>) и возвращает площадь прямоугольника.',
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
// Area вычисляет площадь прямоугольника
func Area(length float64, width float64) float64 {
    // Вставьте ваш код здесь
}
$$,
        '{"function_name": "Area", "parameters": [{"name": "length", "type": "float64"}, {"name": "width", "type": "float64"}]}'::jsonb
    )
),
tests_1_2 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_1_2),
        '[
          { "input": { "length": 5.0, "width": 10.0 }, "output": 50.0 },
          { "input": { "length": 2.5, "width": 4.0 }, "output": 10.0 },
          { "input": { "length": 0.0, "width": 100.0 }, "output": 0.0 }
        ]'::jsonb
    )
),

lesson_2 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Условные конструкции и циклы',
        'Изучение операторов <code>if</code>, <code>else</code>, <code>switch</code> для ветвления и использование единого цикла <code>for</code> в Go.',
        2,
        $$
        <h1>Урок 2: Условные конструкции и циклы</h1>
        <p>Как и любой современный язык, Go поддерживает условные операторы, но имеет более минималистичный подход к циклам.</p>

        <h2>Условный оператор <code>if/else</code></h2>
        <p>В Go условие в <code>if</code> не заключается в скобки (в отличие от C/Java), но фигурные скобки <code>{}</code> обязательны.</p>
        <pre><code>score := 85

if score >= 90 {
    fmt.Println("Отлично")
} else if score >= 70 {
    fmt.Println("Хорошо")
} else {
    fmt.Println("Нужно постараться")
}</code></pre>
        <p><b>Ключевая особенность (короткое объявление):</b> Go позволяет объявить переменную, которая будет видна только в блоке <code>if</code> и <code>else</code>, прямо перед условием.</p>
        <pre><code>if result, err := someFunction(); err != nil {
    // Используем err
} else {
    // Используем result
}</code></pre>

        <h2>Оператор множественного выбора <code>switch</code></h2>
        <p><code>switch</code> в Go мощнее, чем в других языках. По умолчанию, как только одно условие <code>case</code> совпадает, выполнение прекращается (нет "проваливания" — no fallthrough).</p>
        <pre><code>day := "Monday"
switch day {
case "Saturday", "Sunday": // Несколько значений в одном case
    fmt.Println("Выходной")
case "Monday":
    fmt.Println("Начало недели")
default: // Как else
    fmt.Println("Рабочий день")
}</code></pre>

        <h2>Единственный цикл: <code>for</code></h2>
        <p>В Go нет циклов <code>while</code> или <code>do-while</code>. Все формы циклов реализуются с помощью одного ключевого слова <code>for</code>.</p>
        <ol>
            <li><b>Стандартный цикл (как в C):</b>
                <pre><code>for i := 0; i < 5; i++ {
    fmt.Println(i)
}</code></pre>
            </li>
            <li><b>Цикл как <code>while</code>:</b> Просто опускаем инициализацию и пост-операцию.
                <pre><code>sum := 1
for sum < 100 {
    sum += sum
}</code></pre>
            </li>
            <li><b>Бесконечный цикл:</b>
                <pre><code>for {
    // ...
    break // Выход из цикла
}</code></pre>
            </li>
        </ol>
        <p>Операторы <code>break</code> и <code>continue</code> работают так же, как и в других языках.</p>
        $$
    )
    RETURNING id
),
problem_2_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_2),
        'Проверка знака числа',
        'Напишите функцию <code>CheckSign</code>, которая принимает число <code>x</code> типа <code>int</code> и возвращает строку: "Positive", "Negative" или "Zero". Используйте конструкцию <code>if-else if-else</code>.',
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
// CheckSign проверяет знак числа и возвращает соответствующую строку.
func CheckSign(x int) string {
    // Вставьте ваш код здесь
}
$$,
        '{"function_name": "CheckSign", "parameters": [{"name": "x", "type": "int"}]}'::jsonb
    )
),
tests_2_1 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_2_1),
        '[
            { "input": { "x": 10 }, "output": "Positive" },
            { "input": { "x": -5 }, "output": "Negative" },
            { "input": { "x": 0 }, "output": "Zero" }
        ]'::jsonb
    )
),

problem_2_2 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_2),
        'Сумма нечетных до N',
        'Напишите функцию <code>SumOddUpToN</code>, которая вычисляет сумму всех нечетных целых чисел от 1 до <code>n</code> (включительно). Используйте цикл <code>for</code>.',
        'easy',
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
// SumOddUpToN вычисляет сумму нечетных чисел от 1 до n.
func SumOddUpToN(n int) int {
    // Вставьте ваш код здесь
}
$$,
        '{"function_name": "SumOddUpToN", "parameters": [{"name": "n", "type": "int"}]}'::jsonb
    )
),
tests_2_2 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_2_2),
        '[
            { "input": { "n": 5 }, "output": 9 },   
            { "input": { "n": 10 }, "output": 25 }, 
            { "input": { "n": 1 }, "output": 1 },
            { "input": { "n": 0 }, "output": 0 }
        ]'::jsonb
    )
),

lesson_3 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Функции и множественные возвращаемые значения',
        'Изучение синтаксиса объявления функций, передачи аргументов и уникальной особенности Go — возврата нескольких значений.',
        3,
        $$
        <h1>Урок 3: Функции и множественные возвращаемые значения</h1>
        <p>Функции в Go — это основа организации кода. Они должны быть объявлены с типом для каждого параметра и для возвращаемого значения.</p>

        <h2>Объявление функций</h2>
        <p>Синтаксис объявления функции: <code>func functionName(param1 type1, param2 type2) returnType</code>.</p>
        <pre><code>// Функция, не принимающая аргументов и не возвращающая ничего
func logMessage() {
    fmt.Println("Сообщение")
}

// Функция с параметрами и возвращаемым значением
func Multiply(a int, b int) int {
    return a * b
}

// Сокращенный синтаксис, если типы параметров совпадают
func Add(a, b int) int {
    return a + b
}</code></pre>

        <h2>Множественные возвращаемые значения</h2>
        <p>Go позволяет функции возвращать более одного значения. Это часто используется для возврата результата и ошибки одновременно (например, <code>(result, error)</code>).</p>
        <pre><code>func swap(a, b int) (int, int) {
    return b, a
}

x, y := swap(10, 20) // x=20, y=10</code></pre>

        <h2>Именованные возвращаемые значения</h2>
        <p>Вы можете дать имена возвращаемым значениям в сигнатуре функции. В этом случае достаточно использовать пустой оператор <code>return</code> в теле функции.</p>
        <pre><code>func calculate(a, b int) (sum int, diff int) {
    sum = a + b
    diff = a - b
    return // Возвращает sum и diff
}</code></pre>
        $$
    )
    RETURNING id
),
problem_3_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_3),
        'Максимум и минимум',
        'Напишите функцию <code>MinMax</code>, которая принимает два целых числа <code>a</code> и <code>b</code> и возвращает оба: сначала наименьшее, затем наибольшее. Используйте множественные возвращаемые значения.',
        'easy',
        '1 second',
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
// MinMax возвращает наименьшее и наибольшее из двух чисел.
func MinMax(a int, b int) (int, int) {
    // Вставьте ваш код здесь
}
$$,
        '{"function_name": "MinMax", "parameters": [{"name": "a", "type": "int"}, {"name": "b", "type": "int"}]}'::jsonb
    )
),
tests_3_1 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_3_1),
        '[
            { "input": { "a": 10, "b": 5 }, "output": [5, 10] },
            { "input": { "a": 7, "b": 7 }, "output": [7, 7] },
            { "input": { "a": -1, "b": 10 }, "output": [-1, 10] }
        ]'::jsonb
    )
),

problem_3_2 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_3),
        'Разделение на части',
        'Напишите функцию <code>Divide</code>, которая принимает два целых числа <code>numerator</code> и <code>denominator</code>. Функция должна вернуть два значения типа <code>int</code>: результат целочисленного деления и остаток от деления.',
        'easy',
        '1 second',
        67108864
    )
    RETURNING id
),
template_3_2_go AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_3_2),
        'go',
        $$
// Divide возвращает частное и остаток от деления.
func Divide(numerator int, denominator int) (quotient int, remainder int) {
    // Вставьте ваш код здесь
}
$$,
        '{"function_name": "Divide", "parameters": [{"name": "numerator", "type": "int"}, {"name": "denominator", "type": "int"}]}'::jsonb
    )
),
tests_3_2 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_3_2),
        '[
            { "input": { "numerator": 10, "denominator": 3 }, "output": [3, 1] },
            { "input": { "numerator": 15, "denominator": 5 }, "output": [3, 0] },
            { "input": { "numerator": 7, "denominator": 2 }, "output": [3, 1] }
        ]'::jsonb
    )
),

lesson_4 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Массивы и Слайсы (Arrays and Slices)',
        'Изучение разницы между статическими массивами и динамическими слайсами — основным способом работы с коллекциями в Go.',
        4,
        $$
        <h1>Урок 4: Массивы и Слайсы</h1>
        <p>Go имеет две связанные структуры данных для хранения последовательностей: <b>массивы</b> и <b>слайсы</b>. Почти всегда в реальном коде вы будете использовать слайсы.</p>

        <h2>Массивы (Arrays)</h2>
        <p>Массив — это коллекция элементов <b>одного типа</b> с фиксированной длиной. Размер массива является частью его типа, поэтому <code>[3]int</code> и <code>[4]int</code> — это разные типы данных.</p>
        <pre><code>var a [5]int // Массив из 5 нулей (zero value)
b := [3]int{1, 2, 3} // Массив из 3 элементов

fmt.Println(b[1]) // Доступ: 2</code></pre>

        <h2>Слайсы (Slices)</h2>
        <p>Слайс — это мощная и гибкая абстракция над массивами. Слайсы представляют собой ссылки на базовые массивы и могут динамически расти или уменьшаться в размере. Слайс состоит из трех компонентов: указателя на базовый массив, длины (<code>len</code>) и емкости (<code>cap</code>).</p>
        
        <h3>Создание слайсов</h3>
        <ol>
            <li><b>С помощью литерала:</b>
                <pre><code>s := []string{"a", "b", "c"} // Тип []string, без указания размера</code></pre>
            </li>
            <li><b>С помощью <code>make</code>:</b>
                <pre><code>// Создание слайса int с длиной 5 и емкостью 5 (все нули)
s := make([]int, 5) 
// Создание слайса int с длиной 0, но емкостью 10
s := make([]int, 0, 10)</code></pre>
            </li>
        </ol>

        <h3>Добавление элементов: <code>append</code></h3>
        <p>Для добавления элемента используйте встроенную функцию <code>append()</code>. Она может создавать новый базовый массив, если текущая емкость недостаточна.</p>
        <pre><code>numbers := []int{1, 2}
numbers = append(numbers, 3) // numbers: [1 2 3]</code></pre>

        <h3>Срезы (Slicing)</h3>
        <p>Слайс можно создать из другого массива или слайса, используя оператор среза <code>[low:high]</code>. Срез включает элемент с индексом <code>low</code>, но исключает элемент с индексом <code>high</code>.</p>
        <pre><code>s := []int{10, 20, 30, 40, 50}
s1 := s[1:4] // [20 30 40]
s2 := s[:2]  // [10 20]
s3 := s[3:]  // [40 50]</code></pre>

        <h2>Цикл <code>for...range</code> для итерации</h2>
        <p>Для перебора элементов в слайсах и массивах используется цикл <code>for range</code>. Он возвращает индекс и значение для каждой итерации.</p>
        <pre><code>nums := []int{10, 20, 30}
for index, value := range nums {
    fmt.Printf("Индекс: %d, Значение: %d\n", index, value)
}</code></pre>
        <p>Если вам нужен только индекс или только значение, используйте символ подчеркивания <code>_</code> для игнорирования неиспользуемой переменной.</p>
        $$
    )
    RETURNING id
),
problem_4_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_4),
        'Сумма элементов слайса',
        'Напишите функцию <code>SliceSum</code>, которая принимает слайс целых чисел <code>nums</code> (<code>[]int</code>) и возвращает сумму всех его элементов. Используйте цикл <code>for range</code>.',
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
// SliceSum вычисляет сумму всех элементов в слайсе.
func SliceSum(nums []int) int {
    // Вставьте ваш код здесь
}
$$,
        '{"function_name": "SliceSum", "parameters": [{"name": "nums", "type": "[]int"}]}'::jsonb
    )
),
tests_4_1 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_4_1),
        '[
            { "input": { "nums": [1, 2, 3, 4] }, "output": 10 },
            { "input": { "nums": [10, -5, 5] }, "output": 10 },
            { "input": { "nums": [] }, "output": 0 }
        ]'::jsonb
    )
),

problem_4_2 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_4),
        'Фильтрация положительных чисел',
        'Напишите функцию <code>FilterPositive</code>, которая принимает слайс чисел <code>nums</code> (<code>[]int</code>) и возвращает новый слайс, содержащий только положительные числа (больше 0).',
        'medium',
        '1 second',
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
// FilterPositive возвращает новый слайс, содержащий только положительные числа.
func FilterPositive(nums []int) []int {
    // Вставьте ваш код здесь
}
$$,
        '{"function_name": "FilterPositive", "parameters": [{"name": "nums", "type": "[]int"}]}'::jsonb
    )
),
tests_4_2 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_4_2),
        '[
            { "input": { "nums": [1, -2, 3, 0, 5, -1] }, "output": [1, 3, 5] },
            { "input": { "nums": [-10, -20, 0] }, "output": [] },
            { "input": { "nums": [100] }, "output": [100] }
        ]'::jsonb
    )
),

lesson_5 AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES (
        (SELECT id FROM course_insert),
        'Карты (Maps): Коллекции "ключ-значение"',
        'Изучение карт в Go — аналога словарей в Python, включая их создание, доступ и безопасное использование.',
        5,
        $$
        <h1>Урок 5: Карты (Maps)</h1>
        <p>Карта (<code>map</code>) в Go — это неупорядоченная коллекция пар "ключ-значение", известная как хэш-таблица. Карты используются, когда нужно быстро находить значение по его ключу. Ключи должны быть одного типа, значения — другого, но оба типа задаются при создании карты.</p>

        <h2>Создание карт</h2>
        <p>Карты создаются с помощью литерала или встроенной функции <code>make()</code>. Неинициализированная переменная карты имеет нулевое значение <code>nil</code>.</p>
        <pre><code>// 1. С помощью литерала:
countryCapitals := map[string]string{
    "USA": "Washington",
    "France": "Paris",
}

// 2. С помощью make (рекомендуется для инициализации):
userScores := make(map[string]int) // Создает пустую карту</code></pre>

        <h2>Добавление, изменение и доступ</h2>
        <p>Добавление и изменение значений осуществляется простым присваиванием по ключу.</p>
        <pre><code>userScores["Alice"] = 95 // Добавление
userScores["Alice"] = 99 // Изменение

score := userScores["Bob"] // Доступ, score будет 0 (zero value), если Bob нет</code></pre>

        <h2>Проверка существования ключа (Comma Ok Idiom)</h2>
        <p>Прямой доступ к несуществующему ключу вернет нулевое значение, что может быть обманчиво. В Go есть специальный синтаксис для безопасной проверки наличия ключа.</p>
        <pre><code>score, exists := userScores["Bob"]

if exists {
    fmt.Println("У Bob счет:", score)
} else {
    fmt.Println("Bob не найден")
}</code></pre>

        <h2>Удаление элемента</h2>
        <p>Для удаления пары ключ-значение используйте встроенную функцию <code>delete()</code>.</p>
        <pre><code>delete(countryCapitals, "USA")</code></pre>

        <h2>Итерация по карте</h2>
        <p>Используйте цикл <code>for range</code>, который возвращает ключ и значение.</p>
        <pre><code>for country, capital := range countryCapitals {
    fmt.Printf("%s - %s\n", country, capital)
}</code></pre>
        <p>Помните, что порядок обхода элементов в карте не гарантируется.</p>
        $$
    )
    RETURNING id
),
problem_5_1 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_5),
        'Общее количество товаров',
        'Напишите функцию <code>TotalItems</code>, которая принимает карту <code>inventory</code> (<code>map[string]int</code>), где ключ — это название товара, а значение — его количество. Функция должна вернуть общее количество всех товаров на складе.',
        'easy',
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
// TotalItems вычисляет общее количество товаров в инвентаре.
func TotalItems(inventory map[string]int) int {
    // Вставьте ваш код здесь
}
$$,
        '{"function_name": "TotalItems", "parameters": [{"name": "inventory", "type": "map[string]int"}]}'::jsonb
    )
),
tests_5_1 AS (
    INSERT INTO tests (problem_id, tests)
    VALUES (
        (SELECT id FROM problem_5_1),
        '[
            { "input": { "inventory": { "apple": 10, "banana": 20, "orange": 5 } }, "output": 35 },
            { "input": { "inventory": { "pen": 100 } }, "output": 100 },
            { "input": { "inventory": {} }, "output": 0 }
        ]'::jsonb
    )
),

problem_5_2 AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES (
        (SELECT id FROM lesson_5),
        'Проверка наличия пользователя',
        'Напишите функцию <code>HasUser</code>, которая принимает карту <code>users</code> (<code>map[string]int</code>), где ключ — имя пользователя, и строку <code>name</code>. Функция должна вернуть <code>true</code>, если пользователь с таким именем существует в карте, и <code>false</code> в противном случае. Используйте <b>Comma Ok Idiom</b>.',
        'easy',
        '1 second',
        67108864
    )
    RETURNING id
),
template_5_2_go AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES (
        (SELECT id FROM problem_5_2),
        'go',
        $$
// HasUser проверяет, существует ли пользователь с заданным именем в карте.
func HasUser(users map[string]int, name string) bool {
    // Вставьте ваш код здесь
}
$$,
        '{"function_name": "HasUser", "parameters": [{"name": "users", "type": "map[string]int"}, {"name": "name", "type": "string"}]}'::jsonb
    )
)
INSERT INTO tests (problem_id, tests)
VALUES (
    (SELECT id FROM problem_5_2),
    '[
        { "input": { "users": { "Alice": 1, "Bob": 2 }, "name": "Alice" }, "output": true },
        { "input": { "users": { "Alice": 1, "Bob": 2 }, "name": "Charlie" }, "output": false },
        { "input": { "users": {}, "name": "Test" }, "output": false }
    ]'::jsonb
);

COMMIT;
-- +goose StatementEnd
