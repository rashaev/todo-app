package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ListenAddress string `envconfig:"TODO_LISTEN_ADDRESS" default:":8000"`
	Database      DBConfig
	Logging       LogConfig
}

type DBConfig struct {
	Host     string `envconfig:"TODO_DB_HOST" default:":localhost"`
	Port     string `envconfig:"TODO_DB_PORT" default:":5432"`
	Username string `envconfig:"TODO_DB_USERNAME" required:"true"`
	Password string `envconfig:"TODO_DB_PASSWORD" required:"true"`
	DBName   string `envconfig:"TODO_DB_DBNAME" required:"true"`
}

type LogConfig struct {
	Level string `envconfig:"TODO_LOG_LEVEL" default:":info"`
}

func Load() (*Config, error) {
	cfg := new(Config)

	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("error occured while parsing env config: %w", err)
	}
	return cfg, nil
}
