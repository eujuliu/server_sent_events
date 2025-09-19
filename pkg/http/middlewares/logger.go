package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

var Logger gin.HandlerFunc = gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
	return fmt.Sprintf("[%s] %s %s %s %d %s %s\n",
		param.TimeStamp.Format(time.RFC1123),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.ErrorMessage,
	)
})
