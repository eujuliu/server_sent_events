package sse

import (
	"context"
	"encoding/json"
	"sse/internal/interfaces"
	"sync"
)

type SSEService struct {
	consumer      interfaces.IQueue
	clients       map[*Client]bool
	activeClients int
	broadcast     chan interfaces.Event
	register      chan *Client
	unregister    chan *Client
	mu            sync.Mutex
}

func NewSSEService(consumer interfaces.IQueue) *SSEService {
	svc := &SSEService{
		consumer:   consumer,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan interfaces.Event, 100),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	go svc.run()

	return svc
}

func (s *SSEService) run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true
			s.activeClients++
			s.mu.Unlock()
		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				s.activeClients--
				delete(s.clients, client)
				close(client.send)
			}
			s.mu.Unlock()
		case event := <-s.broadcast:
			data, _ := json.Marshal(event)
			s.mu.Lock()
			for client := range s.clients {
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
	return s.consumer.Consume(ctx, queue, func(event interfaces.Event) error {
		select {
		case s.broadcast <- event:
		default:
		}

		return nil
	})
}

func (s *SSEService) RegisterClient(client *Client) {
	s.register <- client
}

func (s *SSEService) UnregisterClient(client *Client) {
	s.unregister <- client
}

func (s *SSEService) ClientsCounts() int {
	return s.activeClients
}
