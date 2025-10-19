package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"problum/internal/config"
	"problum/internal/database"
	"problum/internal/server"

	"github.com/gofiber/fiber/v3"

	"github.com/rs/zerolog/log"
)

type App struct {
	httpServer *fiber.App
	cfg        *config.Config
	db         *database.DB
}

func New() (*App, error) {
	cfg, err := config.New()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create config")
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	db, err := database.New(cfg.DB)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create db")
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	return &App{
		httpServer: server.New(cfg),
		cfg:        cfg,
		db:         db,
	}, nil
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		listenPath := fmt.Sprintf("%s:%d", a.cfg.Server.Host, a.cfg.Server.Port)
		log.Info().Interface("config", a.cfg).Msgf("starting server on: %s", listenPath)
		if err := a.httpServer.Listen(listenPath); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error().Err(err).Msg("failed to listen and serve http server")
			}
		}
	}()

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Info().Msg("gracefully shutdown server")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, a.cfg.Server.ShutdownTimeout)
	defer shutdownCancel()
	if err := a.httpServer.ShutdownWithContext(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown server")
		return err
	}

	log.Info().Msg("successfully stopped server")
	return nil
}
