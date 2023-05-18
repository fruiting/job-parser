package queue

import (
	"fmt"

	"fruiting/job-parser/internal"
	"github.com/mailru/easyjson"
)

type ParseByPositionTaskProducer struct {
	pool internal.RedisPool
}

func NewParseByPositionTaskProducer(pool internal.RedisPool) *ParseByPositionTaskProducer {
	return &ParseByPositionTaskProducer{
		pool: pool,
	}
}

func (p *ParseByPositionTaskProducer) Produce(payload *internal.ParseByPositionTask) error {
	payloadJson, err := easyjson.Marshal(payload)
	if err != nil {
		return fmt.Errorf("can't marshal parse by position task payload: %w", err)
	}

	err = p.pool.Put(string(payloadJson))
	if err != nil {
		return fmt.Errorf("can't put payload into parse by position task queue: %w", err)
	}

	return nil
}
