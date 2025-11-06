package dto

import (
	"time"

	problemDTO "problum/internal/problem/service/dto"

	"problum/internal/api"
	"problum/internal/model"
)

func ToAPI(lesson *Lesson) api.LessonGetResponse {
	return api.LessonGetResponse{
		ID:          lesson.ID,
		CourseID:    lesson.CourseID,
		Name:        lesson.Name,
		Description: lesson.Description,
		Position:    lesson.Position,
		Content:     lesson.Content,
		CreatedAt:   lesson.CreatedAt,
		UpdatedAt:   lesson.UpdatedAt,
		Problems:    problemDTO.ToAPIList(lesson.Problems),
	}
}

func ToAPIList(lessons []*Lesson) []api.LessonGetResponse {
	ans := make([]api.LessonGetResponse, 0, len(lessons))

	for _, lesson := range lessons {
		ans = append(ans, ToAPI(lesson))
	}

	return ans
}

func ToDTO(lesson *model.Lesson, problems []*problemDTO.Problem) *Lesson {
	return &Lesson{
		ID:          lesson.ID,
		CourseID:    lesson.CourseID,
		Name:        lesson.Name,
		Description: lesson.Description,
		Position:    lesson.Position,
		Content:     lesson.Content,
		CreatedAt:   lesson.CreatedAt,
		UpdatedAt:   lesson.UpdatedAt,
		Problems:    problems,
	}
}

func ToDTOList(lessons []*model.Lesson, m map[int][]*problemDTO.Problem) []*Lesson {
	ans := make([]*Lesson, 0, len(lessons))

	for _, lesson := range lessons {
		var problems []*problemDTO.Problem
		if m != nil {
			problems = m[lesson.ID]
		}

		ans = append(ans, ToDTO(lesson, problems))
	}

	return ans
}

type Lesson struct {
	ID          int
	CourseID    int
	Name        string
	Description string
	Position    int
	Content     string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Problems []*problemDTO.Problem
}
