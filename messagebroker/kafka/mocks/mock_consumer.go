package mocks

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type MockConsumer struct {
	ReadMessageFunc func(ctx context.Context) (kafka.Message, error)
}

func (m *MockConsumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	if m.ReadMessageFunc != nil {
		return m.ReadMessageFunc(ctx)
	}
	return kafka.Message{}, nil
}