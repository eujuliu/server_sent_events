package ratelimiter

import (
	"context"
	"fmt"
	"sse/pkg/redis"
	"strconv"
	"time"
)

type SlidingWindowCounterLimiter struct {
	redis         *redis.Redis
	limit         int
	windowSize    int64
	subWindowSize int64
}

func NewSlidingWindowCounterLimiter(
	redis *redis.Redis,
	limit int,
	windowSize, subWindowSize int64,
) *SlidingWindowCounterLimiter {
	return &SlidingWindowCounterLimiter{
		redis:         redis,
		limit:         limit,
		windowSize:    windowSize,
		subWindowSize: subWindowSize,
	}
}

func (rt *SlidingWindowCounterLimiter) GetLimit() int64 {
	return int64(rt.limit)
}

func (rt *SlidingWindowCounterLimiter) Allowed(ctx context.Context, id string) (int64, bool) {
	key := fmt.Sprintf("rate_limit:%s", id)

	subs, err := rt.redis.HGetAll(ctx, key)
	if err != nil {
		panic(err)
	}

	var totalCount int64
	for _, v := range subs {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			panic(err)
		}

		totalCount += n
	}

	allowed := totalCount < int64(rt.limit)

	if allowed {
		currentTime := time.Now().UnixMilli()
		subWindowSizeMillis := rt.subWindowSize * 1000
		currentSubWindow := currentTime / subWindowSizeMillis

		rt.redis.BeginTransaction()

		_, err := rt.redis.HIncrBy(ctx, key, fmt.Sprint(currentSubWindow), 1)
		if err != nil {
			panic(err)
		}

		_, err = rt.redis.HExpire(
			ctx,
			key,
			time.Duration(rt.windowSize)*time.Second,
			"NX",
			fmt.Sprint(currentSubWindow),
		)
		if err != nil {
			panic(err)
		}

		err = rt.redis.ExecTransaction(ctx)
		if err != nil {
			panic(err)
		}
	}

	return int64(rt.limit) - totalCount, allowed
}
