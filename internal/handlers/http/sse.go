package http_handlers

import (
	"io"
	"sse/pkg/sse"

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
	client := sse.NewClient()
	h.sseService.RegisterClient(client)
	defer func() { h.sseService.UnregisterClient(client) }()

	c.Stream(func(w io.Writer) bool {
		if msg, ok := <-client.Send(); ok {
			c.SSEvent("message", msg)
		}

		return true
	})
}
