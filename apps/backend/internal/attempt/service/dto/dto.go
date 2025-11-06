package dto

import (
	"time"

	"problum/internal/api"
	"problum/internal/model"
)

type Attempt struct {
	ID           int
	UserID       int
	ProblemID    int
	Duration     time.Duration
	MemoryUsage  int64
	Language     string
	Code         string
	Status       string
	ErrorMessage *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func ToDTO(attempt *model.Attempt) *Attempt {
	return &Attempt{
		ID:           attempt.ID,
		UserID:       attempt.UserID,
		ProblemID:    attempt.ProblemID,
		Duration:     attempt.Duration,
		MemoryUsage:  attempt.MemoryUsage,
		Language:     attempt.Language,
		Code:         attempt.Code,
		Status:       attempt.Status,
		ErrorMessage: attempt.ErrorMessage,
		CreatedAt:    attempt.CreatedAt,
		UpdatedAt:    attempt.UpdatedAt,
	}
}

func ToDTOList(attempts []*model.Attempt) []*Attempt {
	ans := make([]*Attempt, 0, len(attempts))

	for _, attempt := range attempts {
		ans = append(ans, ToDTO(attempt))
	}

	return ans
}

func ToModel(attempt *Attempt) *model.Attempt {
	return &model.Attempt{
		ID:           attempt.ID,
		UserID:       attempt.UserID,
		ProblemID:    attempt.ProblemID,
		Duration:     attempt.Duration,
		MemoryUsage:  attempt.MemoryUsage,
		Language:     attempt.Language,
		Code:         attempt.Code,
		Status:       attempt.Status,
		ErrorMessage: attempt.ErrorMessage,
		CreatedAt:    attempt.CreatedAt,
		UpdatedAt:    attempt.UpdatedAt,
	}
}

func ToAPI(attempt *Attempt) api.AttemptGetResponse {
	return api.AttemptGetResponse{
		ID:           attempt.ID,
		UserID:       attempt.UserID,
		ProblemID:    attempt.ProblemID,
		Duration:     attempt.Duration,
		MemoryUsage:  attempt.MemoryUsage,
		Language:     attempt.Language,
		Code:         attempt.Code,
		Status:       attempt.Status,
		ErrorMessage: attempt.ErrorMessage,
		CreatedAt:    attempt.CreatedAt,
		UpdatedAt:    attempt.UpdatedAt,
	}
}

func ToAPIList(attempts []*Attempt) []api.AttemptGetResponse {
	ans := make([]api.AttemptGetResponse, 0, len(attempts))

	for _, attempt := range attempts {
		ans = append(ans, ToAPI(attempt))
	}

	return ans
}
