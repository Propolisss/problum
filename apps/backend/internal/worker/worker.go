package worker

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"problum/internal/config"
	"problum/internal/database"
	"problum/internal/nats"
	"problum/internal/redis"
	"problum/internal/solver"

	attemptRepository "problum/internal/attempt/repository"
	attemptService "problum/internal/attempt/service"
	attemptDTO "problum/internal/attempt/service/dto"

	solverDTO "problum/internal/solver/dto"

	templateRepository "problum/internal/template/repository"
	templateService "problum/internal/template/service"

	testRepository "problum/internal/test/repository"
	testService "problum/internal/test/service"

	problemRepository "problum/internal/problem/repository"
	problemService "problum/internal/problem/service"
	problemDTO "problum/internal/problem/service/dto"

	"github.com/bytedance/sonic"
	natsgo "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
)

type AttemptService interface {
	Submit(context.Context, *attemptDTO.Attempt) (int, error)
	Update(ctx context.Context, attempt *attemptDTO.Attempt) error
}

type Solver interface {
	Solve(context.Context, *attemptDTO.Attempt) (*solverDTO.Result, error)
}

type ProblemService interface {
	GetWithOptions(context.Context, int, ...problemService.Option) (*problemDTO.Problem, error)
}

type Worker struct {
	cfg           *config.Config
	db            *database.DB
	rdb           *redis.Redis
	nc            *natsgo.Conn
	js            jetstream.JetStream
	attemptStream jetstream.Stream
	consumer      jetstream.Consumer
	attemptSvc    AttemptService
	problemSvc    ProblemService
	solver        Solver
}

func New() (*Worker, error) {
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

	attemptRepo := attemptRepository.New(db)
	attemptSvc := attemptService.New(attemptRepo)

	templateRepo := templateRepository.New(db)
	templateSvc := templateService.New(templateRepo)

	testRepo := testRepository.New(db)
	testSvc := testService.New(testRepo)

	problemRepo := problemRepository.New(db)
	problemSvc := problemService.New(problemRepo, js, attemptSvc, templateSvc)

	solver := solver.New(testSvc, templateSvc, problemSvc)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := js.Stream(ctx, "ATTEMPTS")
	if err != nil {
		log.Error().Err(err).Msg("Failed to get stream for attempts")
		return nil, fmt.Errorf("failed to get stream for attempts")
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "worker",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to create or update consumer for attempts stream")
		return nil, fmt.Errorf("failed to create or update consumer for attempts stream")
	}

	worker := &Worker{
		cfg:           cfg,
		db:            db,
		rdb:           rdb,
		nc:            nc,
		js:            js,
		attemptSvc:    attemptSvc,
		attemptStream: stream,
		consumer:      consumer,
		solver:        solver,
	}

	return worker, nil
}

func (w *Worker) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		w.work(ctx)
	}()

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	log.Info().Msg("Gracefully shutdown server")
	cancel()

	log.Info().Msg("Successfully stopped worker")
	return nil
}

func (w *Worker) work(ctx context.Context) {
	iter, err := w.consumer.Messages()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create iterator")
		return
	}
	defer iter.Stop()

	go func() {
		<-ctx.Done()
		iter.Stop()
	}()

	log.Info().Msg("Starting pulling messages...")
	for {
		msg, err := iter.Next()
		if err != nil {
			log.Error().Err(err).Msg("Failed to get messages from iterator")
			break
		}

		log.Info().Msg("Pulled message")

		message := &attemptDTO.Attempt{}
		if err := sonic.Unmarshal(msg.Data(), message); err != nil {
			log.Error().Err(err).Msg("Failed to unmarshal message")
		} else {
			log.Info().Interface("message", message).Msg("unmarshaled message")
		}

		result, err := w.solver.Solve(ctx, message)
		if err != nil {
			log.Error().Err(err).Msg("Failed to solve problem")
			continue
		}

		message.Duration = result.Duration
		message.MemoryUsage = result.MemoryUsage
		message.Status = result.Status
		message.ErrorMessage = result.ErrorMessage

		if err := w.attemptSvc.Update(ctx, message); err != nil {
			log.Error().Err(err).Msg("Failed to update attempt")
			continue
		}

		if err := msg.Ack(); err != nil {
			log.Error().Err(err).Msg("Failed to ack message")
		}

		log.Info().Msg("Acked message")
	}
}
