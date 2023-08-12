package main

type Config struct {
	LogLevel   string `long:"log-level" description:"Log level: panic, fatal, warn or warning, info, debug" env:"LOG_LEVEL" required:"true"`
	LogJSON    bool   `long:"log-json" description:"Enable force log format JSON" env:"LOG_JSON"`
	HttpListen string `long:"http-listen" description:"HTTP listen port" env:"HTTP_LISTEN"`

	RedisHost     string `long:"redis-host" description:"Redis host address" env:"REDIS_HOST"`
	RedisPort     string `long:"redis-port" description:"Redis port" env:"REDIS_PORT"`
	RedisUsername string `long:"redis-username" description:"Redis username" env:"REDIS_USERNAME"`
	RedisPassword string `long:"redis-password" description:"Redis password" env:"REDIS_PASSWORD"`

	KafkaBroker          string `long:"kafka-broker" description:"Kafka broker" env:"KAFKA_BROKER"`
	KafkaMaxMessageBytes int    `long:"kafka-max-size-message" description:"Max size message for Kafka" env:"KAFKA_MAX_MESSAGE_BYTES" required:"true"`
	KafkaMaxRetry        int    `long:"kafka-max-retry" description:"Max retry count to connect to Kafka" env:"KAFKA_MAX_RETRY" required:"true"`
	KafkaTopicParseJob   string `long:"kafka-topic-parse-job" description:"Kafka parse job topic" env:"KAFKA_TOPIC_PARSE_JOB" required:"true"`

	HhLinksPerPageBatchSize uint8 `long:"hh-links-per-page-batch-size" description:"HH links per page batch size" env:"HH_LINKS_PER_PAGE_BATCH_SIZE" required:"true"`

	EnablePprof bool `long:"enable-pprof" description:"Enable pprof server" env:"ENABLE_PPROF"`
}
