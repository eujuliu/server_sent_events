package rabbitmq

import (
	"context"
	"encoding/json"
	"sse/internal/config"
	"sse/internal/interfaces"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQ(config *config.RabbitMQConfig) (*RabbitMQ, error) {
	conn, err := amqp.Dial(config.Url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rmq *RabbitMQ) AddDurableQueue(name, exchange, routingKey string) error {
	q, err := rmq.ch.QueueDeclare(strings.ToLower(name), true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = rmq.ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = rmq.ch.QueueBind(q.Name, routingKey, exchange, false, nil)

	return err
}

func (rmq *RabbitMQ) Consume(
	ctx context.Context,
	queue string,
	handler func(interfaces.Event) error,
) error {
	msgs, err := rmq.ch.Consume(queue, "sse", false, false, false, false, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case msg := <-msgs:
			var event interfaces.Event

			if err := json.Unmarshal(msg.Body, &event); err != nil {
				_ = msg.Nack(false, true)
				continue
			}

			if err := handler(event); err != nil {
				_ = msg.Nack(false, true)
				continue
			}

			_ = msg.Ack(false)
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
