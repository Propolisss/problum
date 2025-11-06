package dto

import (
	"time"

	"problum/internal/api"
	"problum/internal/model"
	templateDTO "problum/internal/template/service/dto"
)

type ProblemSubmit struct {
	ID        int
	ProblemID int
	UserID    int
	Language  string
	Code      string
}

type Problem struct {
	ID          int
	LessonID    int
	Name        string
	Statement   string
	Difficulty  string
	TimeLimit   time.Duration
	MemoryLimit int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Template    *templateDTO.Template
}

func ToDTO(problem *model.Problem, template *templateDTO.Template) *Problem {
	return &Problem{
		ID:          problem.ID,
		LessonID:    problem.LessonID,
		Name:        problem.Name,
		Statement:   problem.Statement,
		Difficulty:  problem.Difficulty,
		TimeLimit:   problem.TimeLimit,
		MemoryLimit: problem.MemoryLimit,
		CreatedAt:   problem.CreatedAt,
		UpdatedAt:   problem.UpdatedAt,
		Template:    template,
	}
}

func ToDTOList(problems []*model.Problem) []*Problem {
	ans := make([]*Problem, 0, len(problems))

	for _, problem := range problems {
		ans = append(ans, ToDTO(problem, nil))
	}

	return ans
}

func ToAPI(problem *Problem) api.ProblemGetResponse {
	return api.ProblemGetResponse{
		ID:          problem.ID,
		LessonID:    problem.LessonID,
		Name:        problem.Name,
		Statement:   problem.Statement,
		Difficulty:  problem.Difficulty,
		TimeLimit:   problem.TimeLimit,
		MemoryLimit: problem.MemoryLimit,
		CreatedAt:   problem.CreatedAt,
		UpdatedAt:   problem.UpdatedAt,
		Template:    templateDTO.ToAPI(problem.Template),
	}
}

func ToAPIList(problems []*Problem) []api.ProblemGetResponse {
	ans := make([]api.ProblemGetResponse, 0, len(problems))

	for _, problem := range problems {
		ans = append(ans, ToAPI(problem))
	}

	return ans
}
