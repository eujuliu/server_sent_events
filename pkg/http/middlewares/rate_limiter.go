package middlewares

import (
	"fmt"
	"net/http"
	ratelimiter "sse/pkg/rate_limiter"

	"github.com/gin-gonic/gin"
)

func RateLimiter(limiter *ratelimiter.SlidingWindowCounterLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		count, allowed := limiter.Allowed(c, ip)

		c.Header("X-Ratelimit-Remaining", fmt.Sprint(max(count-1, 0)))
		c.Header("X-Ratelimit-Limit", fmt.Sprint(limiter.GetLimit()))

		if allowed {
			c.Next()
			return
		}

		c.JSON(http.StatusTooManyRequests, gin.H{
			"error": "Rate Limit Exceeded",
		})
		c.Abort()
	}
}
