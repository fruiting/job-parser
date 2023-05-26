package main

// Config application configuration
type Config struct {
	LogLevel   string `long:"log-level" description:"Log level: panic, fatal, warn or warning, info, debug" env:"LOG_LEVEL" required:"true"`
	LogJSON    bool   `long:"log-json" description:"Enable force log format JSON" env:"LOG_JSON"`
	HttpListen string `long:"http-listen" description:"HTTP listen port" env:"HTTP_LISTEN"`

	RedisHost     string `long:"redis-host" description:"Redis host address" env:"REDIS_HOST"`
	RedisPort     string `long:"redis-port" description:"Redis port" env:"REDIS_PORT"`
	RedisUsername string `long:"redis-username" description:"Redis username" env:"REDIS_USERNAME"`
	RedisPassword string `long:"redis-password" description:"Redis password" env:"REDIS_PASSWORD"`

	PgDbHost     string `long:"pg-db-host" description:"PG DB host address" env:"PG_DB_HOST"`
	PgDbPort     int    `long:"pg-db-port" description:"PG DB port" env:"PG_DB_PORT"`
	PgDbUsername string `long:"pg-db-username" description:"PG DB username" env:"PG_DB_USERNAME"`
	PgDbPassword string `long:"pg-db-password" description:"PG DB password" env:"PG_DB_PASSWORD"`
	PgDbName     string `long:"pg-db-name" description:"PG DB name" env:"PG_DB_NAME"`

	TgApiKey string `long:"tg-api-key" description:"Telegram api key" env:"TG_API_KEY"`

	EnablePprof bool `long:"enable-pprof" description:"Enable pprof server" env:"ENABLE_PPROF"`
}
