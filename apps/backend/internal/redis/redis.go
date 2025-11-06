package redis

import (
	"context"
	"fmt"
	"time"

	"problum/internal/config"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Redis struct {
	cfg *config.Redis
	rdb *redis.Client
}

func New(cfg *config.Redis) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		log.Error().Err(err).Msg("Failed to ping redis")
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return &Redis{
		cfg: cfg,
		rdb: rdb,
	}, nil
}

func (r *Redis) Ping(ctx context.Context) error {
	return r.rdb.Ping(ctx).Err()
}

func (r *Redis) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return r.rdb.Set(ctx, key, value, ttl).Err()
}

func (r *Redis) SetString(ctx context.Context, key, value string, ttl time.Duration) error {
	return r.Set(ctx, key, value, ttl)
}

func (r *Redis) Get(ctx context.Context, key string) ([]byte, error) {
	return r.rdb.Get(ctx, key).Bytes()
}

func (r *Redis) GetString(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

func (r *Redis) Delete(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}

func (r *Redis) SAdd(ctx context.Context, key string, members ...interface{}) error {
	return r.rdb.SAdd(ctx, key, members).Err()
}

func (r *Redis) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.rdb.SMembers(ctx, key).Result()
}

func (r *Redis) SRem(ctx context.Context, key string, members ...interface{}) error {
	return r.rdb.SRem(ctx, key, members...).Err()
}
