package internal

import (
	"context"
	"fmt"

	"github.com/adjust/redismq"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
)

type Processor struct {
	consumer       Consumer
	storage        Storage
	chatBotHandler ChatBotHandler
	priceSorter    PriceSorter
	parser         *GeneralParser
	logger         *zap.Logger
}

func NewProcessor(
	consumer Consumer,
	storage Storage,
	chatBotHandler ChatBotHandler,
	priceSorter PriceSorter,
	parser *GeneralParser,
	logger *zap.Logger,
) *Processor {
	return &Processor{
		consumer:       consumer,
		storage:        storage,
		chatBotHandler: chatBotHandler,
		priceSorter:    priceSorter,
		parser:         parser,
		logger:         logger,
	}
}

func (p *Processor) Run(ctx context.Context) error {
	for {
		task, err := p.consumer.Get()
		if err != nil {
			continue
		}
		if task == nil {
			continue
		}

		if task.Payload == "" {
			p.logger.Warn("empty payload")
			continue
		}

		payload := &Payload{}
		err = easyjson.Unmarshal([]byte(task.Payload), payload)
		if err != nil {
			p.logger.Error("can't unmarshal payload", zap.Error(err))
			err = p.requeue(task)
			if err != nil {
				p.logger.Error("can't requeue task", zap.Error(err))
			}

			continue
		}

		ctxLogger := p.logger.With(zap.Any("payload", payload))

		go func() {
			jobs, err := p.parser.Run(payload.PositionName)
			if err != nil {
				ctxLogger.Error("can't parse", zap.Error(err))
				err = p.requeue(task)
				if err != nil {
					ctxLogger.Error("can't requeue task", zap.Error(err))
				}

				return
			}

			min, max, median := p.priceSorter.PricesFromJobs(jobs)
			jobsInfo := &JobsInfo{
				PositionToParse: payload.PositionName,
				MinSalary:       min,
				MaxSalary:       max,
				MedianSalary:    median,
				PopularSkills:   p.popularSkills(jobs),
				Parser:          p.parser.parser.Parser(),
				Jobs:            jobs,
			}

			err = p.storage.Set(ctx, jobsInfo)
			if err != nil {
				ctxLogger.Error("can't set into storage", zap.Error(err))
				err = p.requeue(task)
				if err != nil {
					ctxLogger.Error("can't requeue task", zap.Error(err))
				}

				return
			}

			err = p.chatBotHandler.Push(jobsInfo)
			if err != nil {
				err = p.requeue(task)
				if err != nil {
					ctxLogger.Error("can't requeue task", zap.Error(err))
				}

				ctxLogger.Error("can't push jobs info into chat bot handler", zap.Error(err))
			}

			err = task.Ack()
			if err != nil {
				p.logger.Error("can't ack task", zap.Error(err))
			}

			ctxLogger.Info("task completed")
		}()
	}
}

func (p *Processor) popularSkills(jobs []*Job) []string {
	return nil
}

func (p *Processor) requeue(task *redismq.Package) error {
	err := task.Requeue()
	if err != nil {
		err = task.Fail()
		if err != nil {
			return fmt.Errorf("can't fail task: %w", err)
		}

		return fmt.Errorf("can't requeue task: %w", err)
	}

	return nil
}
