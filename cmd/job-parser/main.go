package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
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
	"github.com/jmoiron/sqlx"
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

	pool := redismq.CreateQueue(
		cfg.RedisHost,
		cfg.RedisPort,
		cfg.RedisPassword,
		0,
		"parse_by_position_tasks_queue",
	)
	parseByPositionTaskProducer := queue.NewParseByPositionTaskProducer(pool)
	consumer, err := pool.AddConsumer("parse_by_position_tasks_consumer")
	if err != nil {
		logger.Fatal("can't add consumer for redis queue", zap.Error(err))
	}

	pgDb, err := initPgDb(cfg.PgDbHost, cfg.PgDbPort, cfg.PgDbUsername, cfg.PgDbPassword, cfg.PgDbName)
	if err != nil {
		//logger.Fatal("can't init pg db", zap.Error(err))
	}

	pgsqlStorage := pgsql.NewStorage(pgDb, logger)
	httpClient := &http.Client{}
	chatBotHandler := telegram.NewChatBotHandler(
		cfg.HttpListen,
		httpClient,
		cfg.TgApiKey,
		parseByPositionTaskProducer,
		pgsqlStorage,
		logger,
	)
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

	queueConsumer := internal.NewQueueConsumer([]internal.Consumer{
		parseByPositionConsumer,
	})

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		queueConsumer.Run(ctx)
	}()

	wg.Add(1)
	go func() {
		//http.ListenAndServe(":8080", http.HandlerFunc(Handler))
		logger.Info("Starting chat bot handle server")
		err := chatBotHandler.ListenAndServe()
		cancelFunc() // stop app if handle server was stopped
		if err != nil {
			logger.Error("Error on listen and serve chat bot handle server", zap.Error(err))
		}
	}()

	wg.Wait()
}

func initPgDb(host string, port int, user, password, dbName string) (*sqlx.DB, error) {
	hostPort := fmt.Sprintf("%s:%d", host, port)
	db := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(user, password),
		Host:     hostPort,
		RawQuery: fmt.Sprintf("database=%s", dbName),
	}
	dbConnection, err := sqlx.Open("pgx", db.String())
	if err != nil {
		return nil, fmt.Errorf("can't open pgsql db connection: %w", err)
	}

	dbConnection.SetConnMaxLifetime(6 * time.Minute)
	dbConnection.SetConnMaxIdleTime(4 * time.Minute)
	dbConnection.SetMaxOpenConns(50)

	return dbConnection, nil
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
