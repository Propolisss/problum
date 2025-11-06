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
	"problum/internal/nats"
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

	courseHandler "problum/internal/course/delivery/http"
	courseRepository "problum/internal/course/repository"
	courseService "problum/internal/course/service"

	lessonHandler "problum/internal/lesson/delivery/http"
	lessonRepository "problum/internal/lesson/repository"
	lessonService "problum/internal/lesson/service"
	lessonDTO "problum/internal/lesson/service/dto"

	enrollmentHandler "problum/internal/enrollment/delivery/http"
	enrollmentRepository "problum/internal/enrollment/repository"
	enrollmentService "problum/internal/enrollment/service"
	enrollmentDTO "problum/internal/enrollment/service/dto"

	problemHandler "problum/internal/problem/delivery/http"
	problemRepository "problum/internal/problem/repository"
	problemService "problum/internal/problem/service"
	problemDTO "problum/internal/problem/service/dto"

	attemptHandler "problum/internal/attempt/delivery/http"
	attemptRepository "problum/internal/attempt/repository"
	attemptService "problum/internal/attempt/service"
	attemptDTO "problum/internal/attempt/service/dto"

	templateRepository "problum/internal/template/repository"
	templateService "problum/internal/template/service"

	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
)

type AttemptSvc interface {
	Get(context.Context, int) (*attemptDTO.Attempt, error)
}

type EnrollmentService interface {
	Get(ctx context.Context, courseID, userID int) (*enrollmentDTO.Enrollment, error)
}
type LessonService interface {
	Get(ctx context.Context, id int) (*lessonDTO.Lesson, error)
}
type ProblemService interface {
	Get(context.Context, int) (*problemDTO.Problem, error)
	Submit(context.Context, *problemDTO.ProblemSubmit) (int, error)
}

type App struct {
	httpServer *fiber.App
	cfg        *config.Config
	db         *database.DB
	rdb        *redis.Redis
	nc         *natsgo.Conn
	js         jetstream.JetStream
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

	nc, err := nats.New(cfg.Nats)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create nats")
		return nil, fmt.Errorf("failed to create nats: %w", err)
	}

	js, err := nats.NewStream(nc)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create jetstream")
		return nil, fmt.Errorf("failed to create jetstream: %w", err)
	}

	sessionRepo := sessionRepository.New(db)
	sessionSvc := sessionService.New(sessionRepo)

	userRepo := userRepository.New(db)
	userSvc := userService.New(userRepo)

	authSvc := authService.New(rdb, userSvc, sessionSvc)
	authHdl := authHandler.New(cfg, authSvc)

	attemptRepo := attemptRepository.New(db)
	attemptSvc := attemptService.New(attemptRepo)
	attemptHdl := attemptHandler.New(cfg, attemptSvc)

	templateRepo := templateRepository.New(db)
	templateSvc := templateService.New(templateRepo)

	problemRepo := problemRepository.New(db)
	problemSvc := problemService.New(problemRepo, js, attemptSvc, templateSvc)
	problemHdl := problemHandler.New(cfg, problemSvc)

	lessonRepo := lessonRepository.New(db)
	lessonSvc := lessonService.New(lessonRepo, problemSvc)
	lessonHdl := lessonHandler.New(cfg, lessonSvc)

	enrollmentRepo := enrollmentRepository.New(db)
	enrollmentSvc := enrollmentService.New(enrollmentRepo)
	enrollmentHdl := enrollmentHandler.New(cfg, enrollmentSvc)

	courseRepo := courseRepository.New(db)
	courseSvc := courseService.New(courseRepo, lessonSvc, enrollmentSvc)
	courseHdl := courseHandler.New(cfg, courseSvc)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     "ATTEMPTS",
		Subjects: []string{"ATTEMPTS.*"},
	}); err != nil {
		log.Error().Err(err).Msg("Failed to create or update stream")
		return nil, fmt.Errorf("failed to create or update stream: %w", err)
	}

	app := &App{
		httpServer: server.New(cfg),
		cfg:        cfg,
		db:         db,
		rdb:        rdb,
		nc:         nc,
		js:         js,
	}

	setupRoutes(
		app,
		authHdl,
		courseHdl,
		lessonHdl,
		enrollmentHdl,
		problemHdl,
		enrollmentSvc,
		lessonSvc,
		problemSvc,
		attemptHdl,
		attemptSvc,
	)

	return app, nil
}

func setupRoutes(
	app *App,
	authHdl *authHandler.Handler,
	courseHdl *courseHandler.Handler,
	lessonHdl *lessonHandler.Handler,
	enrollmentHdl *enrollmentHandler.Handler,
	problemHdl *problemHandler.Handler,
	enrollmentSvc EnrollmentService,
	lessonSvc LessonService,
	problemSvc ProblemService,
	attemptHdl *attemptHandler.Handler,
	attemptSvc AttemptSvc,
) {
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

	// auth
	auth := app.httpServer.Group("/auth")
	auth.Post("/login", authHdl.Login)
	auth.Post("/refresh", authHdl.Refresh)
	auth.Post("/logout", middleware.Auth(app.rdb), authHdl.Logout)
	auth.Post("/register", authHdl.Register)

	// course
	course := app.httpServer.Group("/courses")
	course.Use(middleware.Auth(app.rdb))
	course.Get("/", courseHdl.List)
	course.Get("/:courseID", middleware.Course(enrollmentSvc), courseHdl.Get)

	// lesson
	lesson := course.Group("/:courseID/lessons")
	lesson.Use(middleware.Course(enrollmentSvc))
	lesson.Get("/:lessonID", middleware.Lesson(lessonSvc), lessonHdl.Get)

	// problem
	problem := course.Group("/:courseID/problems")
	problem.Use(middleware.Course(enrollmentSvc))
	problem.Get("/:problemID", middleware.Problem(problemSvc, lessonSvc), problemHdl.Get)
	problem.Post("/:problemID/submit", middleware.Problem(problemSvc, lessonSvc), problemHdl.Submit)

	// attempt
	attempt := app.httpServer.Group("/attempts")
	attempt.Use(middleware.Auth(app.rdb))
	attempt.Get("/", attemptHdl.ListByUserID)
	attempt.Get("/:attemptID", middleware.Attempt(attemptSvc), attemptHdl.Get)
	problem.Get("/:problemID", middleware.Problem(problemSvc, lessonSvc), problemHdl.Get)
	problem.Get("/:problemID/attempts", middleware.Problem(problemSvc, lessonSvc), attemptHdl.ListByProblemID)
	problem.Post("/:problemID/submit", middleware.Problem(problemSvc, lessonSvc), problemHdl.Submit)

	// enrollment
	enrollment := app.httpServer.Group("/enrollments")
	enrollment.Use(middleware.Auth(app.rdb))
	enrollment.Post("/", enrollmentHdl.Enroll)
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
