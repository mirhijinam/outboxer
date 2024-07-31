package config

import (
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig           DBConfig
	ServerConfig       ServerConfig
	LoggerConfig       LoggerConfig
	EventHandlerConfig EventHandlerConfig
	KafkaConfig        KafkaConfig
}

type DBConfig struct {
	PgUser     string `env:"PGUSER"`
	PgPassword string `env:"PGPASSWORD"`
	PgHost     string `env:"PGHOST"`
	PgPort     uint16 `env:"PGPORT"`
	PgDatabase string `env:"PGDATABASE"`
	PgSSLMode  string `env:"PGSSLMODE"`
}

type ServerConfig struct {
	Port        string `env:"HTTP_PORT" envDefault:"8080"`
	Timeout     string `env:"TIMEOUT" envDefault:"5s"`
	IdleTimeout string `env:"IDLE_TIMEOUT" envDefault:"30s"`
}

type LoggerConfig struct {
	Mode     string `env:"LOG_MODE" envDefault:"info"`
	Filepath string `env:"LOG_FILE"`
}

type EventHandlerConfig struct {
	CooldownSec int `env:"EVENT_HANDLER_CD_SEC"`
}

type KafkaConfig struct {
	Brokers     string        `env:"KAFKA_BROKERS"`
	Topic       string        `env:"KAFKA_TOPIC"`
	GroupID     string        `env:"KAFKA_GROUP_ID"`
	Timeout     time.Duration `env:"KAFKA_TIMEOUT"`
	OffsetReset string        `env:"KAFKA_OFFSET_RESET"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func New() (*Config, error) {
	config := &Config{}

	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("failed to parse config from environment variables: %w", err)
	}

	return config, nil
}
