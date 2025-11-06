package dto

import (
	"time"

	"problum/internal/model"
)

func ToDTO(enrollment *model.Enrollment) *Enrollment {
	return &Enrollment{
		ID:        enrollment.ID,
		CourseID:  enrollment.CourseID,
		UserID:    enrollment.UserID,
		CreatedAt: enrollment.CreatedAt,
		UpdatedAt: enrollment.UpdatedAt,
	}
}

func ToDTOList(enrollments []*model.Enrollment) []*Enrollment {
	ans := make([]*Enrollment, 0, len(enrollments))

	for _, enrollment := range enrollments {
		ans = append(ans, ToDTO(enrollment))
	}

	return ans
}

type Enrollment struct {
	ID        int
	CourseID  int
	UserID    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
