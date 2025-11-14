package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	LogLevel        slog.Level    `env:"APP_LOG_LEVEL"`
	Address         string        `env:"APP_ADDRESS,required"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT,required"`
	DB              *DBConfig
}

type DBConfig struct {
	Host string `env:"DB_HOST,required"`
	Port int    `env:"DB_PORT,required"`
	User string `env:"DB_USER,required"`
	Pass string `env:"DB_PASSWORD,required"`
	Name string `env:"DB_NAME,required"`
}

func Load() (Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse base config: %w", err)
	}

	var dbCfg DBConfig
	if err := env.Parse(&dbCfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse DB config: %w", err)
	}
	cfg.DB = &dbCfg

	return cfg, nil
}
