package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	ReadMessage(ctx context.Context) (kafka.Message, error)
}
