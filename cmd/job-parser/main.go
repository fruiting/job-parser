package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"fruiting/job-parser/internal"
	"fruiting/job-parser/internal/chatbothandler/telegram"
	"fruiting/job-parser/internal/parser/headhunter"
	"fruiting/job-parser/internal/pricesorter"
	"fruiting/job-parser/internal/queue"
	"fruiting/job-parser/internal/skillssorter"
	"fruiting/job-parser/internal/storage/pgsql"
	"github.com/adjust/redismq"
	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	var cfg Config
	parser := flags.NewParser(&cfg, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		log.Fatal(fatalJsonLog("Failed to parse config", err))
	}

	logger, err := initLogger(cfg.LogLevel, cfg.LogJSON)
	if err != nil {
		log.Fatal(fatalJsonLog("Failed to init logger", err))
	}

	pool := redismq.CreateQueue(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, 0, "parse_by_position_tasks_queue2")
	consumer, err := pool.AddConsumer("parse_by_position_tasks_consumer2")
	if err != nil {
		logger.Fatal("can't add consumer for redis queue", zap.Error(err))
	}

	pgsqlStorage := pgsql.NewStorage()
	chatBotHandler := telegram.NewChatBotHandle(logger)
	headHunterParser := headhunter.NewParser(logger)
	generalParser := internal.NewGeneralParser(headHunterParser)

	parseByPositionConsumer := queue.NewParseByPositionTaskConsumer(
		consumer,
		pgsqlStorage,
		chatBotHandler,
		pricesorter.NewPriceSorter(),
		skillssorter.NewSkillsSorter(),
		generalParser,
		logger,
	)

	payload := "asdas"
	payloadJson, err := json.Marshal(payload)
	err = pool.Put(string(payloadJson))

	queueConsumer := internal.NewQueueConsumer([]internal.Consumer{
		parseByPositionConsumer,
	})

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		queueConsumer.Run(context.Background())
	}()

	wg.Wait()
}

// initLogger создает и настраивает новый экземпляр логгера
func initLogger(logLevel string, isLogJson bool) (*zap.Logger, error) {
	lvl := zap.InfoLevel
	err := lvl.UnmarshalText([]byte(logLevel))
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal log-level: %w", err)
	}
	opts := zap.NewProductionConfig()
	opts.Level = zap.NewAtomicLevelAt(lvl)
	opts.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if opts.InitialFields == nil {
		opts.InitialFields = map[string]interface{}{}
	}
	if !isLogJson {
		opts.Encoding = "console"
		opts.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	return opts.Build()
}

func fatalJsonLog(msg string, err error) string {
	escape := func(s string) string {
		return strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), `"`, `\"`)
	}
	errString := ""
	if err != nil {
		errString = err.Error()
	}
	return fmt.Sprintf(
		`{"level":"fatal","ts":"%s","msg":"%s","error":"%s"}`,
		time.Now().Format(time.RFC3339),
		escape(msg),
		escape(errString),
	)
}
