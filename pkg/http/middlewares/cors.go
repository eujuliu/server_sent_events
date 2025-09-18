package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Cors gin.HandlerFunc = cors.New(cors.Config{
	AllowOrigins:     []string{"http://localhost:5500"},
	AllowMethods:     []string{"GET", "OPTIONS"},
	AllowHeaders:     []string{"Origin", "Content-Type"},
	ExposeHeaders:    []string{"Content-Length"},
	AllowCredentials: true,
	MaxAge:           12 * time.Hour,
})
