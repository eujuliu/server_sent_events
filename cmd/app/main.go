package main

import (
	"context"
	"log/slog"
	"os"
	"sse/internal/config"
	http_handlers "sse/internal/handlers/http"
	"sse/pkg/http"
	"sse/pkg/http/middlewares"
	"sse/pkg/rabbitmq"
	ratelimiter "sse/pkg/rate_limiter"
	"sse/pkg/redis"
	"sse/pkg/sse"

	"github.com/gin-gonic/gin"
)

func main() {
	config := config.NewConfig()

	loggerLevel := slog.LevelInfo

	if config.Server.GinMode == gin.DebugMode {
		loggerLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: loggerLevel,
	}))

	slog.SetDefault(logger)

	rmq, err := rabbitmq.NewRabbitMQ(config.RabbitMQ)
	if err != nil {
		panic(err)
	}

	err = rmq.AddDurableQueue("events", "events", "events.send")
	if err != nil {
		panic(err)
	}

	rdb, err := redis.NewRedis(config.Redis)
	if err != nil {
		panic(err)
	}

	sseService := sse.NewSSEService(rmq, rdb)

	go func() {
		err = sseService.Start(context.Background(), "events")
		if err != nil {
			panic(err)
		}
	}()

	sseHandler := http_handlers.NewSSEHandler(sseService)

	limiter := ratelimiter.NewSlidingWindowCounterLimiter(
		rdb,
		config.RateLimiter.RequestLimit,
		config.RateLimiter.WindowSize,
		config.RateLimiter.SubWindowSize,
	)

	server := http.NewServer(config.Server, limiter)

	server.Router().
		GET("/events", middlewares.SSEHeaders, middlewares.Authentication, sseHandler.Handle)

	server.Listen()
}
