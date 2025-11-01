package config

import (
	"sse/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RabbitMQConfig struct {
	Port     string
	User     string
	Password string
	Host     string
	Url      string
}

type ServerConfig struct {
	Host    string
	Port    string
	GinMode string
}

type RedisConfig struct {
	Addr     string
	Password string
	Username string
	DB       int
}

type RateLimiterConfig struct {
	RequestLimit  int
	WindowSize    int64
	SubWindowSize int64
}

type Config struct {
	RabbitMQ    *RabbitMQConfig
	Server      *ServerConfig
	Redis       *RedisConfig
	RateLimiter *RateLimiterConfig
}

func NewConfig() *Config {
	redis_db, err := strconv.Atoi(utils.GetEnv("REDIS_DB", "0"))
	if err != nil {
		panic(err)
	}

	return &Config{
		RabbitMQ: &RabbitMQConfig{
			Port:     utils.GetEnv("RABBITMQ_PORT", "5672"),
			User:     utils.GetEnv("RABBITMQ_DEFAULT_USER", "local_user"),
			Password: utils.GetEnv("RABBITMQ_DEFAULT_PASS", "local_password"),
			Host:     utils.GetEnv("RABBITMQ_HOST", "localhost"),
			Url: utils.GetEnv(
				"RABBITMQ_CONNECTION_STRING",
				"amqp://guest:guest@localhost:5672/",
			),
		},
		Server: &ServerConfig{
			Host:    utils.GetEnv("HOST", "0.0.0.0"),
			Port:    utils.GetEnv("PORT", "8080"),
			GinMode: utils.GetEnv("GIN_MODE", gin.ReleaseMode),
		},
		Redis: &RedisConfig{
			Addr:     utils.GetEnv("REDIS_ADDRESS", "localhost:6379"),
			Password: utils.GetEnv("REDIS_PASSWORD", ""),
			Username: utils.GetEnv("REDIS_USERNAME", ""),
			DB:       redis_db,
		},
		RateLimiter: &RateLimiterConfig{
			RequestLimit:  5,
			WindowSize:    60, // in seconds
			SubWindowSize: 20, // in seconds
		},
	}
}
