package main

// Config конфигурация приложения
type Config struct {
	LogLevel string `long:"log-level" description:"Log level: panic, fatal, warn or warning, info, debug" env:"LOG_LEVEL" required:"true"`
	LogJSON  bool   `long:"log-json" description:"Enable force log format JSON" env:"LOG_JSON"`

	RedisHost     string `long:"redis-host" description:"Redis host address" env:"REDIS_HOST"`
	RedisPort     string `long:"redis-port" description:"Redis port" env:"REDIS_PORT"`
	RedisUsername string `long:"redis-username" description:"Redis username" env:"REDIS_USERNAME"`
	RedisPassword string `long:"redis-password" description:"Redis password" env:"REDIS_PASSWORD"`
}
