package queue

import (
	"context"
	"encoding/json"
	"sse/internal/interfaces"
	"sync"
)

type InMemoryQueue struct {
	queues map[string]chan []byte
	mu     sync.RWMutex
	done   chan struct{}
}

func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{queues: make(map[string]chan []byte), done: make(chan struct{})}
}

func (q *InMemoryQueue) ensureQueue(name string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, ok := q.queues[name]; !ok {
		q.queues[name] = make(chan []byte, 100)
	}

	go func() {
		<-q.done
		close(q.queues[name])
	}()
}

func (q *InMemoryQueue) Publish(key string, exchangeName string, data []byte) error {
	q.ensureQueue(key)
	q.queues[key] <- data
	return nil
}

func (q *InMemoryQueue) Consume(
	ctx context.Context,
	queue string,
	handler func(interfaces.Event) error,
) error {
	q.ensureQueue(queue)
	msgs := q.queues[queue]

	for {
		select {
		case msg := <-msgs:
			var event interfaces.Event

			if err := json.Unmarshal(msg, &event); err != nil {
				continue
			}

			if err := handler(event); err != nil {
				continue
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (q *InMemoryQueue) Close() {
	close(q.done)
}
