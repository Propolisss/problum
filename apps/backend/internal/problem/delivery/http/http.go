package http

import (
	"context"
	"fmt"
	"strconv"

	"problum/internal/api"
	"problum/internal/config"
	"problum/internal/problem/service/dto"

	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
)

type Service interface {
	GetWithTemplate(context.Context, int, string) (*dto.Problem, error)
	Submit(context.Context, *dto.ProblemSubmit) (int, error)
}

type Handler struct {
	cfg *config.Config
	svc Service
}

func New(cfg *config.Config, svc Service) *Handler {
	return &Handler{
		cfg: cfg,
		svc: svc,
	}
}

func (h *Handler) Get(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("problemID"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse problem id in query param")
		return fmt.Errorf("failed to parse problem id in query param")
	}

	language := c.Query("language")

	problem, err := h.svc.GetWithTemplate(c.Context(), id, language)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get problem")
		return fmt.Errorf("failed to get problem")
	}

	return c.JSON(dto.ToAPI(problem))
}

func (h *Handler) Submit(c fiber.Ctx) error {
	submitReq := &api.ProblemSubmitRequest{}
	if err := c.Bind().JSON(submitReq); err != nil {
		return err
	}

	userID, ok := c.Locals("user_id").(int)
	if !ok {
		log.Error().Msg("Failed to get userID")
		return fmt.Errorf("Failed to get userID")
	}

	problemID, err := strconv.Atoi(c.Params("problemID"))
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse problem id in query param")
		return fmt.Errorf("failed to parse problem id in query param")
	}

	attemptID, err := h.svc.Submit(c.Context(), &dto.ProblemSubmit{
		ProblemID: problemID,
		UserID:    userID,
		Language:  submitReq.Language,
		Code:      submitReq.Code,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to submit")
		return fmt.Errorf("failed to submit")
	}

	return c.JSON(api.ProblemSubmitResponse{
		AttemptID: attemptID,
	})
}
