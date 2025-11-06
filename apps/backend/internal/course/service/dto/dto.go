package dto

import (
	"time"

	"problum/internal/api"
	enrollmentDTO "problum/internal/enrollment/service/dto"
	lessonDTO "problum/internal/lesson/service/dto"
	"problum/internal/model"
)

func ToDTO(course *model.Course, lessons []*lessonDTO.Lesson, enrolled bool) *CourseDTO {
	return &CourseDTO{
		ID:          course.ID,
		Name:        course.Name,
		Description: course.Description,
		Tags:        course.Tags,
		Status:      course.Status,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
		Lessons:     lessons,
		Enrolled:    enrolled,
	}
}

func ToDTOList(courses []*model.Course, enrollments []*enrollmentDTO.Enrollment) []*CourseDTO {
	ans := make([]*CourseDTO, 0, len(courses))

	me := make(map[int]bool, len(enrollments))
	for _, enrollment := range enrollments {
		me[enrollment.CourseID] = true
	}

	for _, course := range courses {
		ans = append(ans, ToDTO(course, nil, me[course.ID]))
	}

	return ans
}

func ToAPI(course *CourseDTO) api.CourseGetResponse {
	return api.CourseGetResponse{
		ID:          course.ID,
		Name:        course.Name,
		Description: course.Description,
		Tags:        course.Tags,
		Status:      course.Status,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
		Lessons:     lessonDTO.ToAPIList(course.Lessons),
		Enrolled:    course.Enrolled,
	}
}

func ListToAPI(courses []*CourseDTO) []api.CourseGetResponse {
	ans := make([]api.CourseGetResponse, 0, len(courses))

	for _, course := range courses {
		ans = append(ans, ToAPI(course))
	}

	return ans
}

type CourseDTO struct {
	ID          int
	Name        string
	Description string
	Tags        []string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Enrolled    bool

	Lessons []*lessonDTO.Lesson
}
