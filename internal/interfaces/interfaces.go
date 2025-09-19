package interfaces

import "context"

type Event struct {
	ClientID string `json:"clientId"`
	Type     string `json:"type"`
	Data     string `json:"data"`
}

type IQueue interface {
	Consume(ctx context.Context, queue string, handler func(Event) error) error
}
