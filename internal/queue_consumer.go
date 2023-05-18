package internal

import (
	"context"

	"github.com/adjust/redismq"
)

//go:generate mockgen -source=queue_consumer.go -destination=./queue_consumer_mock.go -package=internal
//go:generate easyjson -output_filename=./queue_easyjson.go

type RedisPool interface {
	Put(payload string) error
}

// RedisConsumer queue consumer by redis
type RedisConsumer interface {
	Get() (*redismq.Package, error)
	GetUnacked() (*redismq.Package, error)
}

type Consumer interface {
	Consume(ctx context.Context)
}

//easyjson:json ParseByPositionTask
type ParseByPositionTask struct {
	PositionName string `json:"position_name"`
}

type QueueConsumer struct {
	consumers []Consumer
}

func NewQueueConsumer(consumers []Consumer) *QueueConsumer {
	return &QueueConsumer{
		consumers: consumers,
	}
}

func (p *QueueConsumer) Run(ctx context.Context) {
	for {
		for _, consumer := range p.consumers {
			go consumer.Consume(ctx)
		}
	}
}
