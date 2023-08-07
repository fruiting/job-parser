package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"fruiting/job-parser/internal"
	httpinternal "fruiting/job-parser/internal/api/http"
	"fruiting/job-parser/internal/parser/hh"
	"fruiting/job-parser/internal/parser/indeed"
	"fruiting/job-parser/internal/queue"
	"fruiting/job-parser/internal/queue/kafka"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

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

	httpClient := &http.Client{}
	httpClientInternal := httpinternal.NewClient(httpClient)

	hhParser := hh.NewParser()
	_ = indeed.NewParser()

	hhParsingProcessor := internal.NewParsingProcessor(hhParser, httpClientInternal)
	//hhParsingProcessor.Run("golang developer")

	kafkaConsumer, err := initKafkaConsumer(
		cfg.KafkaBroker,
		cfg.KafkaMaxRetry,
		cfg.KafkaMaxMessageBytes,
		[]string{cfg.KafkaTopicParseJob},
		[]*internal.ParsingProcessor{hhParsingProcessor},
		logger,
	)
	if err != nil {
		logger.Fatal("can't init kafka consumer")
	}

	httpServer := httpinternal.NewServer(cfg.HttpListen, cfg.EnablePprof, logger)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		logger.Info("Starting http server", zap.String("port", cfg.HttpListen))
		err := httpServer.ListenAndServe()
		cancelFunc() // stop app if handle server was stopped
		if err != nil {
			logger.Error("Error on listen and serve http server", zap.Error(err))
		}
	}()

	wg.Add(1)
	go func() {
		logger.Info("Starting mq consumer", zap.String("port", cfg.HttpListen))
		err := kafkaConsumer.Run(ctx)
		if err != nil {
			logger.Error("can't consume", zap.Error(err))
		}
	}()

	wg.Wait()
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

func initKafkaConsumer(
	broker string,
	maxRetry int,
	maxMessageBytes int,
	topics []string,
	parsingProcessors []*internal.ParsingProcessor,
	logger *zap.Logger,
) (*kafka.Consumer, error) {
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Producer.Retry.Max = maxRetry
	kafkaCfg.Producer.RequiredAcks = sarama.WaitForAll
	kafkaCfg.Producer.Return.Successes = true
	kafkaCfg.Producer.MaxMessageBytes = maxMessageBytes

	cg, err := sarama.NewConsumerGroup([]string{broker}, "groupId", kafkaCfg)
	if err != nil {
		return nil, fmt.Errorf("failed init kafka consumer: %w", err)
	}

	return kafka.NewConsumer(cg, topics, queue.NewConsumerHandler(parsingProcessors), logger), nil
}
