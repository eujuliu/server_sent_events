package redis

import (
	"context"
	"fmt"
	"reflect"
	"sse/internal/config"
	"sync"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type Redis struct {
	db     *redis.Client
	config *config.RedisConfig
	tx     redis.Pipeliner
	mu     sync.Mutex
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

func (r *Redis) BeginTransaction() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tx = r.db.TxPipeline()
}

func (r *Redis) ExecTransaction(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.tx == nil {
		return fmt.Errorf("you need to initialize the transaction first")
	}

	_, err := r.tx.Exec(ctx)
	r.tx = nil

	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) DiscardTransaction() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.tx == nil {
		return fmt.Errorf("you need to initialize the transaction first")
	}

	r.tx.Discard()
	r.tx = nil

	return nil
}

func (r *Redis) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.db.HGetAll(ctx, key).Result()
}

func (r *Redis) HIncrBy(ctx context.Context, key, field string, incr int64) (int64, error) {
	var tx redis.Cmdable = r.db

	if r.tx != nil {
		tx = r.tx
	}

	return tx.HIncrBy(ctx, key, field, incr).Result()
}

func (r *Redis) HExpire(
	ctx context.Context,
	key string,
	expiration time.Duration,
	mode string,
	fields ...string,
) ([]int64, error) {
	var tx redis.Cmdable = r.db

	if r.tx != nil {
		tx = r.tx
	}

	args := redis.HExpireArgs{}

	reflect.ValueOf(&args).Elem().FieldByName(mode).SetBool(true)

	return tx.HExpireWithArgs(ctx, key, expiration, args, fields...).Result()
}
