package api

import "github.com/gofiber/fiber/v3"

type EnrollRequest struct {
	CourseID int `json:"course_id"`
}
type EnrollResponse struct{}

type EnrollmentAPI interface {
	Enroll(fiber.Ctx) error
}
