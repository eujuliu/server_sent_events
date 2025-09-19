package sse

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sse/internal/interfaces"
	"sse/pkg/redis"
	"strings"
	"sync"
)

type SSEService struct {
	queue      interfaces.IQueue
	redis      *redis.Redis
	clients    map[string]*Client
	broadcast  chan interfaces.Event
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

func NewSSEService(queue interfaces.IQueue, redis *redis.Redis) *SSEService {
	svc := &SSEService{
		queue:      queue,
		redis:      redis,
		clients:    make(map[string]*Client),
		broadcast:  make(chan interfaces.Event, 100),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	svc.subscribeKeySpaceNotifications()
	go svc.run()

	return svc
}

func (s *SSEService) run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client.id] = client
			s.mu.Unlock()
		case client := <-s.unregister:
			s.mu.Lock()

			if _, ok := s.clients[client.id]; ok {
				close(client.send)
				close(client.close)
				delete(s.clients, client.id)
			}

			s.mu.Unlock()
		case event := <-s.broadcast:
			data, _ := json.Marshal(event)
			s.mu.Lock()

			if client, ok := s.clients[event.ClientID]; ok {
				select {
				case client.send <- data:
				default:
				}
			}

			s.mu.Unlock()
		}
	}
}

func (s *SSEService) Start(ctx context.Context, queue string) error {
	return s.queue.Consume(ctx, queue, func(event interfaces.Event) error {
		select {
		case s.broadcast <- event:
		default:
		}

		return nil
	})
}

func (s *SSEService) RegisterClient(ctx context.Context, client *Client) error {
	_, err := s.redis.Get(ctx, fmt.Sprintf("session_id:%v", client.id))
	if err != nil {
		return err
	}

	s.register <- client

	return nil
}

func (s *SSEService) UnregisterClient(client *Client) {
	s.unregister <- client
}

func (s *SSEService) getClient(id string) (*Client, bool) {
	client, ok := s.clients[id]

	return client, ok
}

func (s *SSEService) subscribeKeySpaceNotifications() {
	pubsub := s.redis.Subscribe(context.Background(), "__keyevent@0__:expired")

	go func() {
		for msg := range pubsub.Channel() {
			if !strings.HasPrefix(msg.Payload, "session_id:") {
				continue
			}

			expiredKey := strings.Split(msg.Payload, ":")[1]
			slog.Info(fmt.Sprintf("session %s closed!", expiredKey))

			if client, ok := s.getClient(expiredKey); ok {
				s.UnregisterClient(client)
			}
		}
	}()
}
