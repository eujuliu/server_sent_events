package redis

import (
	"context"
	"sse/internal/config"

	redis "github.com/redis/go-redis/v9"
)

type Redis struct {
	db     *redis.Client
	config *config.RedisConfig
}

func NewRedis(config *config.RedisConfig) (*Redis, error) {
	db := redis.NewClient(&redis.Options{
		Addr:       config.Addr,
		Password:   config.Password,
		Username:   config.Username,
		DB:         config.DB,
		ClientName: "sse",
	})

	err := db.ConfigSet(context.Background(), "notify-keyspace-events", "KEA").Err()
	if err != nil {
		return nil, err
	}

	return &Redis{
		db:     db,
		config: config,
	}, nil
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.db.Get(ctx, key).Result()
}

func (r *Redis) Subscribe(ctx context.Context, key string) *redis.PubSub {
	return r.db.Subscribe(ctx, key)
}
