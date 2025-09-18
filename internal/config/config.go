package config

import (
	"sse/pkg/utils"

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

type Config struct {
	RabbitMQ *RabbitMQConfig
	Server   *ServerConfig
}

func NewConfig() *Config {
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
			Port:    utils.GetEnv("PORT", "8081"),
			GinMode: utils.GetEnv("GIN_MODE", gin.ReleaseMode),
		},
	}
}
