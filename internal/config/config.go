package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig     DBConfig
	ServerConfig ServerConfig
	LoggerConfig LoggerConfig
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
	Mode     string `env:"LOG_MODE" envDefault:"DEBUG"`
	Filepath string `env:"LOG_FILE"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func New() (*Config, error) {
	dbconf, err := GetDBConfig()
	if err != nil {
		return &Config{}, err
	}

	srvconf, err := GetServerConfig()
	if err != nil {
		return &Config{}, err
	}

	logconf, err := GetLoggerConfig()
	if err != nil {
		return &Config{}, err
	}

	return &Config{
		DBConfig:     dbconf,
		ServerConfig: srvconf,
		LoggerConfig: logconf,
	}, nil
}

func GetDBConfig() (DBConfig, error) {
	pgPort, err := strconv.ParseInt(os.Getenv("PGPORT"), 0, 16)
	if err != nil {
		return DBConfig{}, err
	}

	return DBConfig{
		PgUser:     os.Getenv("PGUSER"),
		PgPassword: os.Getenv("PGPASSWORD"),
		PgHost:     os.Getenv("PGHOST"),
		PgPort:     uint16(pgPort),
		PgDatabase: os.Getenv("PGDATABASE"),
		PgSSLMode:  os.Getenv("PGSSLMODE"),
	}, nil
}

func GetServerConfig() (ServerConfig, error) {
	return ServerConfig{
		Port:        ":" + os.Getenv("HTTP_PORT"),
		Timeout:     os.Getenv("TIMEOUT"),
		IdleTimeout: os.Getenv("IDLE_TIMEOUT"),
	}, nil
}

func GetLoggerConfig() (LoggerConfig, error) {
	return LoggerConfig{
		Mode:     os.Getenv("LOG_MODE"),
		Filepath: os.Getenv("LOG_FILE"),
	}, nil
}
