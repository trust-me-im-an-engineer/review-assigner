package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	LogLevel        slog.Level    `env:"LOG_LEVEL"`
	ServerAddress   string        `env:"SERVER_ADDRESS,required"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT,required"`
	AdminToken      string        `env:"ADMIN_TOKEN,required"`
	DB              DBConfig      `envPrefix:"DB_"`
}

type DBConfig struct {
	Host string `env:"HOST,required"`
	Port int    `env:"PORT,required"`
	User string `env:"USER,required"`
	Pass string `env:"PASSWORD,required"`
	Name string `env:"NAME,required"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse base config: %w", err)
	}
	return cfg, nil
}
