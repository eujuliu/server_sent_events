package middlewares

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var Cors gin.HandlerFunc = cors.New(cors.Config{
	AllowAllOrigins:  true,
	AllowMethods:     []string{"GET"},
	AllowHeaders:     []string{"Origin", "Content-Type"},
	ExposeHeaders:    []string{"Content-Length"},
	AllowCredentials: true,
	MaxAge:           12 * time.Hour,
})
