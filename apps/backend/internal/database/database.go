package database

import (
	"context"
	"fmt"
	"time"

	"problum/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
	"github.com/rs/zerolog/log"
)

type DB struct {
	Pool *pgxpool.Pool
	cfg  *config.DB
}

func New(cfg *config.DB) (*DB, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.GetDSN())
	if err != nil {
		log.Error().Err(err).Msg("Failed to create pgx pool config")
		return nil, fmt.Errorf("failed to create pgx pool config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.ConnMaxLifetime

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create pgx pool")
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		log.Error().Err(err).Msg("Failed to ping db")
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return &DB{
		Pool: pool,
		cfg:  cfg,
	}, nil
}

func (db *DB) Migrate() error {
	sqlDB := stdlib.OpenDB(*db.Pool.Config().ConnConfig)
	defer sqlDB.Close()

	if err := goose.SetDialect("postgres"); err != nil {
		log.Error().Err(err).Msg("Failed to set dialect")
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	if err := goose.Up(sqlDB, "./migrations"); err != nil {
		log.Error().Err(err).Msg("Failed to run migrations db")
		return fmt.Errorf("failed to run migrations db: %w", err)
	}

	return nil
}

func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
	}
}

func (db *DB) Ready() error {
	if db.Pool == nil {
		log.Error().Msg("db pool is nil")
		return fmt.Errorf("db pool is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.Pool.Ping(ctx)
}

func (db *DB) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		db.cfg.User,
		db.cfg.Password,
		db.cfg.Host,
		db.cfg.Port,
		db.cfg.DBName,
		db.cfg.SSLMode,
	)
}
