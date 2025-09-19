package http_handlers

import (
	"context"
	"io"
	"net/http"
	"sse/pkg/sse"
	"time"

	"github.com/gin-gonic/gin"
)

type SSEHandler struct {
	sseService *sse.SSEService
}

func NewSSEHandler(sseService *sse.SSEService) *SSEHandler {
	return &SSEHandler{
		sseService: sseService,
	}
}

func (h *SSEHandler) Handle(c *gin.Context) {
	clientId := c.Query("userId")

	client := sse.NewClient(clientId)
	err := h.sseService.RegisterClient(context.Background(), client)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":  http.StatusBadRequest,
			"error": "this client is not logged",
		})

		return
	}

	defer func() { h.sseService.UnregisterClient(client) }()
	pingTicker := time.Tick(15 * time.Second)

	c.Stream(func(w io.Writer) bool {
		select {
		case msg := <-client.Send():
			c.SSEvent("message", msg)
		case <-pingTicker:
			c.SSEvent("ping", "pong")
		case <-client.Close():
			c.SSEvent("disconnected", "true")
			return false
		}

		return true
	})
}
