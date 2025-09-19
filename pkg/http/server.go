package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"sse/internal/config"
	"sse/pkg/http/middlewares"
	ratelimiter "sse/pkg/rate_limiter"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
	server *http.Server
}

func NewServer(
	config *config.ServerConfig,
	limiter *ratelimiter.SlidingWindowCounterLimiter,
) *Server {
	gin.SetMode(config.GinMode)

	router := gin.New()

	router.Use(middlewares.Cors)
	router.Use(middlewares.Logger)
	router.Use(middlewares.RateLimiter(limiter))
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	server := http.Server{
		Addr:    fmt.Sprintf("%s:%s", config.Host, config.Port),
		Handler: router,
	}

	return &Server{
		router: router,
		server: &server,
	}
}

func (s *Server) Listen() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error(fmt.Sprintf("server listen: %s", err))
			panic(err)
		}
	}()

	<-ctx.Done()

	stop()

	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		slog.Error(fmt.Sprintf("server forced to shutdown: %s", err))
		panic(err)
	}

	slog.Info("server exiting")
}

func (s *Server) Router() *gin.RouterGroup {
	return s.router.Group("/")
}
