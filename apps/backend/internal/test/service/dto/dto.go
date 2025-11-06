package dto

import (
	"encoding/json"
	"time"

	"problum/internal/model"

	"github.com/bytedance/sonic"
)

type TestCase struct {
	Input  json.RawMessage `json:"input"`
	Output json.RawMessage `json:"output"`
}

type Test struct {
	ID        int
	ProblemID int
	Tests     []TestCase `json:"tests"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ToDTO(test *model.Test) *Test {
	tests := make([]TestCase, 0)
	sonic.Unmarshal(test.Tests, &tests)

	return &Test{
		ID:        test.ID,
		ProblemID: test.ProblemID,
		Tests:     tests,
		CreatedAt: test.CreatedAt,
		UpdatedAt: test.UpdatedAt,
	}
}
