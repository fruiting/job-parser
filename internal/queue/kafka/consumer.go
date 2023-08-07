package kafka

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

type Consumer struct {
	consumer sarama.ConsumerGroup
	topics   []string
	handler  sarama.ConsumerGroupHandler
	logger   *zap.Logger
}

func NewConsumer(
	consumer sarama.ConsumerGroup,
	topics []string,
	handler sarama.ConsumerGroupHandler,
	logger *zap.Logger,
) *Consumer {
	return &Consumer{
		consumer: consumer,
		topics:   topics,
		handler:  handler,
		logger:   logger,
	}
}

func (c *Consumer) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		err := c.consumer.Consume(ctx, c.topics, c.handler)
		if err != nil {
			return fmt.Errorf("can't consume: %w", err)
		}
	}
}
