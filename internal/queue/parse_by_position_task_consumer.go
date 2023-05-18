package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"fruiting/job-parser/internal"
	"github.com/adjust/redismq"
	"go.uber.org/zap"
)

type ParseByPositionTaskConsumer struct {
	consumer       internal.RedisConsumer
	storage        internal.Storage
	chatBotHandler internal.ChatBotHandler
	priceSorter    internal.PriceSorter
	skillsSorter   internal.SkillsSorter
	generalParser  *internal.GeneralParser
	logger         *zap.Logger
}

func NewParseByPositionTaskConsumer(
	consumer internal.RedisConsumer,
	storage internal.Storage,
	chatBotHandler internal.ChatBotHandler,
	priceSorter internal.PriceSorter,
	skillsSorter internal.SkillsSorter,
	generalParser *internal.GeneralParser,
	logger *zap.Logger,
) *ParseByPositionTaskConsumer {
	return &ParseByPositionTaskConsumer{
		consumer:       consumer,
		storage:        storage,
		chatBotHandler: chatBotHandler,
		priceSorter:    priceSorter,
		skillsSorter:   skillsSorter,
		generalParser:  generalParser,
		logger:         logger,
	}
}

func (c *ParseByPositionTaskConsumer) Consume(ctx context.Context) {
	for {
		task, err := c.consumer.Get()
		if err != nil {
			continue
		}
		if task == nil {
			continue
		}

		if task.Payload == "" {
			c.logger.Warn("empty payload")
			continue
		}

		var payload *internal.ParseByPositionTask
		err = json.Unmarshal([]byte(task.Payload), payload)
		if err != nil {
			c.logger.Error("can't unmarshal payload", zap.Error(err))

			continue
		}

		ctxLogger := c.logger.With(zap.Any("payload", payload))

		go func() {
			err = c.execute(ctx, payload)
			if err != nil {
				err = c.requeue(task)
				if err != nil {
					ctxLogger.Error("can't requeue task", zap.Error(err))
				}
			}

			err = task.Ack()
			if err != nil {
				ctxLogger.Error("can't ack task", zap.Error(err))
			}

			ctxLogger.Info("task completed")
		}()
	}
}

func (c *ParseByPositionTaskConsumer) execute(ctx context.Context, payload *internal.ParseByPositionTask) error {
	jobs, err := c.generalParser.Run(internal.Name(payload.PositionName))
	if err != nil {
		return fmt.Errorf("can't parse: %w", err)
	}

	min, max, median := c.priceSorter.PricesFromJobs(jobs)
	jobsInfo := &internal.JobsInfo{
		PositionToParse: internal.Name(payload.PositionName),
		MinSalary:       min,
		MaxSalary:       max,
		MedianSalary:    median,
		PopularSkills:   c.skillsSorter.MostPopularSkills(jobs, internal.MostPopularSkillsCount),
		Parser:          c.generalParser.Parser.Parser(),
		Jobs:            jobs,
	}

	err = c.storage.Set(
		ctx,
		jobsInfo.PositionToParse,
		jobsInfo.MinSalary,
		jobsInfo.MaxSalary,
		jobsInfo.MedianSalary,
		jobsInfo.Parser,
	)
	if err != nil {
		return fmt.Errorf("can't set into storage: %w", err)
	}

	err = c.chatBotHandler.Push(jobsInfo)
	if err != nil {
		return fmt.Errorf("can't push jobs info into chat bot handler: %w", err)
	}

	return nil
}

func (c *ParseByPositionTaskConsumer) requeue(task *redismq.Package) error {
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
