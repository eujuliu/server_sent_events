package interfaces

import "context"

type Event struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

type IQueue interface {
	Consume(ctx context.Context, queue string, handler func(Event) error) error
}
