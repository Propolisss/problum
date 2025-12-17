import time
import json
import requests
import psycopg2
import sys
import os

API_URL = os.getenv("API_URL", "http://backend:8080")
DB_HOST = os.getenv("DB_HOST", "postgres")
DB_PORT = os.getenv("DB_PORT", "5432")
DB_USER = os.getenv("DB_USER", "problum")
DB_PASS = os.getenv("DB_PASS", "problum")
DB_NAME = os.getenv("DB_NAME", "problum")

DB_DSN = f"postgresql://{DB_USER}:{DB_PASS}@{DB_HOST}:{DB_PORT}/{DB_NAME}?sslmode=disable"

TEST_USER = {"login": "test_whitebox", "password": "password", "repeated_password": "password"}
COURSE_NAME = "WhiteBox Test Course"
LESSON_NAME = "Lesson 1"
PROBLEM_NAME = "Sum Two"

GO_HELPERS = """
var memoryHolder []byte

func EatMemory(nMB int) {
	memoryHolder = make([]byte, nMB*1024*1024)
	for i := range memoryHolder {
		memoryHolder[i] = 1
	}
}

func ActiveSleep(ms int) {
	start := time.Now()
	target := time.Duration(ms) * time.Millisecond
	for time.Since(start) < target {
	}
}
"""

SQL_INIT = """
WITH course_ins AS (
    INSERT INTO courses (name, description, tags, status) 
    VALUES (%s, 'Test Desc', '{"test"}', 'published') RETURNING id
),
lesson_ins AS (
    INSERT INTO lessons (course_id, name, description, position, content)
    VALUES ((SELECT id FROM course_ins), %s, 'Desc', 1, 'Content') RETURNING id
),
problem_ins AS (
    INSERT INTO problems (lesson_id, name, statement, difficulty, time_limit, memory_limit)
    VALUES ((SELECT id FROM lesson_ins), %s, 'Sum a+b', 'easy', '2 seconds', 67108864) RETURNING id
),
template_ins AS (
    INSERT INTO templates (problem_id, language, code, metadata)
    VALUES ((SELECT id FROM problem_ins), 'go', 
    '// Sum func\\nfunc Sum(a, b int) int {\\n return 0 \\n}', 
    %s)
),
test_ins AS (
    INSERT INTO tests (problem_id, tests)
    VALUES ((SELECT id FROM problem_ins), %s)
)
SELECT id FROM problem_ins;
"""

TEMPLATE_META = json.dumps({
    "function_name": "Sum",
    "parameters": [{"name": "a", "type": "int"}, {"name": "b", "type": "int"}],
    "return_type": "int"
})

TEST_CASES = json.dumps([
    {"input": {"a": 1, "b": 2}, "output": 3},
    {"input": {"a": -1, "b": 1}, "output": 0}
])

SCENARIOS = [
    {
        "name": "Positive (AC)", 
        "exp": "AC", 
        "code": "func Sum(a, b int) int { return a + b }"
    },
    {
        "name": "Wrong Answer (WA)", 
        "exp": "WA", 
        "code": "func Sum(a, b int) int { return a + b + 1 }"
    },
    {
        "name": "Compilation Err (CE)", 
        "exp": "CE", 
        "code": "func Sum(a, b int) int { return a + }"
    },
    {
        "name": "Runtime Err (RE)", 
        "exp": "RE", 
        "code": 'func Sum(a, b int) int { panic("panic"); return a + b }'
    },
    {
        "name": "Memory Limit (MLE)", 
        "exp": "MLE", 
        "code": "func Sum(a, b int) int { EatMemory(100); return a+b }"
    },
    {
        "name": "Time Limit (TO)", 
        "exp": ["TO", "RE"], 
        "code": "func Sum(a, b int) int { ActiveSleep(3000); return a+b }"
    }
]

def get_conn():
    for _ in range(15):
        try:
            return psycopg2.connect(DB_DSN)
        except:
            print("Waiting for DB...")
            time.sleep(1)
    raise Exception("DB Unreachable")

def wait_api():
    for _ in range(15):
        try:
            requests.get(f"{API_URL}/readyz", timeout=1)
            return
        except:
            print("Waiting for API...")
            time.sleep(1)
    raise Exception("API Unreachable")

def run():
    wait_api()
    conn = get_conn()
    cur = conn.cursor()
    
    print(">>> Создаём тестовые данные...")
    cur.execute(SQL_INIT, (COURSE_NAME, LESSON_NAME, PROBLEM_NAME, TEMPLATE_META, TEST_CASES))
    problem_id = cur.fetchone()[0]
    conn.commit()

    cur.execute("SELECT course_id FROM lessons JOIN problems ON lessons.id = problems.lesson_id WHERE problems.id = %s", (problem_id,))
    course_id = cur.fetchone()[0]
    cur.close()
    conn.close()

    try:
        requests.post(f"{API_URL}/auth/register", json=TEST_USER)
        resp = requests.post(f"{API_URL}/auth/login", json={"login": TEST_USER['login'], "password": TEST_USER['password']})
        if resp.status_code != 200:
            raise Exception("Login failed")
        token = resp.json()["access_token"]
        headers = {"Authorization": f"Bearer {token}"}

        requests.post(f"{API_URL}/enrollments", json={"course_id": course_id}, headers=headers)

        failed = False
        print("\n>>> Запускаем тесты... <<<\n")
        
        for s in SCENARIOS:
            print(f"Test: {s['name']} -> ", end="")
            sys.stdout.flush()
            
            full_code = GO_HELPERS + "\n" + s["code"]
            
            sub = requests.post(
                f"{API_URL}/courses/{course_id}/problems/{problem_id}/submit", 
                json={
                    "language": "go",
                    "code": full_code,
                },
                headers=headers,
            )
            
            if sub.status_code != 200:
                print(f"Submit Failed: {sub.text}")
                failed = True
                continue

            attempt_id = sub.json()["attempt_id"]
            
            status = "pending"
            for _ in range(30):
                res = requests.get(f"{API_URL}/attempts/{attempt_id}", headers=headers).json()
                if res["status"] != "pending":
                    status = res["status"]
                    break
                time.sleep(0.5)
            
            expected = s["exp"] if isinstance(s["exp"], list) else [s["exp"]]
            
            if status in expected:
                print(f"PASS [{status}]")
            else:
                print(f"FAIL. Got {status}, expected {expected}")
                failed = True

        if failed:
            sys.exit(1)
        print("\n>>> SUCCESS <<<")

    finally:
        conn = get_conn()
        cur = conn.cursor()
        cur.execute("DELETE FROM courses WHERE name = %s", (COURSE_NAME,))
        cur.execute("DELETE FROM users WHERE login = %s", (TEST_USER['login'],))
        conn.commit()
        cur.close()
        conn.close()

if __name__ == "__main__":
    run()
