package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"problum/internal/config"
	"problum/internal/database"
	"problum/internal/middleware"
	"problum/internal/redis"
	"problum/internal/server"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"

	authHandler "problum/internal/auth/delivery/http"
	authService "problum/internal/auth/service"

	userRepository "problum/internal/user/repository"
	userService "problum/internal/user/service"

	sessionRepository "problum/internal/session/repository"
	sessionService "problum/internal/session/service"

	"github.com/rs/zerolog/log"
)

type App struct {
	httpServer *fiber.App
	cfg        *config.Config
	db         *database.DB
	rdb        *redis.Redis
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

	rdb, err := redis.New(cfg.Redis)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create redis")
		return nil, fmt.Errorf("failed to create redis: %w", err)
	}

	sessionRepo := sessionRepository.New(db)
	sessionSvc := sessionService.New(sessionRepo)

	userRepo := userRepository.New(db)
	userSvc := userService.New(userRepo)

	authSvc := authService.New(db, rdb, userSvc, sessionSvc)
	authHdl := authHandler.New(cfg, authSvc)

	app := &App{
		httpServer: server.New(cfg),
		cfg:        cfg,
		db:         db,
		rdb:        rdb,
	}

	setupRoutes(app, authHdl)

	return app, nil
}

func setupRoutes(app *App, authHdl *authHandler.Handler) {
	// healthchecks
	app.httpServer.Get(healthcheck.LivenessEndpoint, healthcheck.New())
	app.httpServer.Get(healthcheck.ReadinessEndpoint, healthcheck.New(healthcheck.Config{
		Probe: func(c fiber.Ctx) bool {
			ctx, cancel := context.WithTimeout(c.Context(), 10*time.Second)
			defer cancel()

			return app.db.Pool.Ping(ctx) == nil && app.rdb.Ping(ctx) == nil
		},
	}))
	app.httpServer.Get(healthcheck.StartupEndpoint, healthcheck.New())

	// temp
	app.httpServer.Get("/secure", middleware.Auth(app.rdb), healthcheck.New())

	// auth
	auth := app.httpServer.Group("/auth")
	auth.Post("/login", authHdl.Login)
	auth.Post("/refresh", authHdl.Refresh)
	auth.Post("/logout", middleware.Auth(app.rdb), authHdl.Logout)
	auth.Post("/register", authHdl.Register)
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := a.db.Migrate(); err != nil {
		log.Error().Err(err).Msg("Failed to run migrations")
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	go func() {
		listenPath := fmt.Sprintf("%s:%d", a.cfg.Server.Host, a.cfg.Server.Port)
		log.Info().Interface("config", a.cfg).Msgf("starting server on: %s", listenPath)
		if err := a.httpServer.Listen(listenPath); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error().Err(err).Msg("Failed to listen and serve http server")
			}
		}
	}()

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Info().Msg("Gracefully shutdown server")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, a.cfg.Server.ShutdownTimeout)
	defer shutdownCancel()
	if err := a.httpServer.ShutdownWithContext(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown server")
		return err
	}

	log.Info().Msg("Successfully stopped server")
	return nil
}
