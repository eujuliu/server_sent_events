package sse_test

import (
	"context"
	"encoding/json"
	"sse/internal/interfaces"
	"sse/internal/queue"
	"sse/pkg/redis"
	"sse/pkg/sse"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"

	. "sse/test"
)

func TestSSE(t *testing.T) {
	rdb := &redis.Redis{}
	db, mock := redismock.NewClientMock()

	ctx, stop := context.WithTimeout(context.Background(), 15*time.Second)
	defer stop()

	err := SetPrivateField(rdb, "db", db)
	Ok(t, err)

	queue := queue.NewInMemoryQueue()

	sseService := sse.NewSSEService(queue, rdb)

	go func() {
		err = sseService.Start(ctx, "events")
		Ok(t, err)
	}()

	mock.ClearExpect()

	client := sse.NewClient("123")

	mock.ExpectGet("session_id:123").SetVal("123")

	err = sseService.RegisterClient(ctx, client)
	Ok(t, err)

	event := interfaces.Event{
		ClientID: "123",
		Type:     "Error",
		Data:     "Hello, World!",
	}

	message, err := json.Marshal(event)

	Ok(t, err)

	err = queue.Publish("events", "", message)

	Ok(t, err)

	var decoded interfaces.Event
	msg := <-client.Send()
	err = json.Unmarshal(msg, &decoded)
	Ok(t, err)

	Equals(t, event, decoded)
	mock.MatchExpectationsInOrder(true)
}
