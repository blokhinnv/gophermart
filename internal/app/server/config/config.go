package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	RunAddress                string        `env:"RUN_ADDRESS"`
	DatabaseURI               string        `env:"DATABASE_URI"`
	AccrualSystemAddress      string        `env:"ACCRUAL_SYSTEM_ADDRESS"`
	JWTSigningKey             string        `env:"JWT_SIGNING_KEY"              envDefault:"practicum"`
	JWTExpireDuration         time.Duration `env:"JWT_EXPIRE_DURATION"          envDefault:"1h"`
	AccrualSystemPoolInterval time.Duration `env:"ACCRUAL_SYSTEM_POOL_INTERVAL" envDefault:"1s"`
}

func NewConfig() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}
	if cfg.RunAddress == "" {
		flag.StringVar(&cfg.RunAddress, "a", "", "адрес и порт запуска сервиса")
	}
	if cfg.DatabaseURI == "" {
		flag.StringVar(&cfg.DatabaseURI, "d", "", "адрес подключения к базе данных")
	}
	if cfg.AccrualSystemAddress == "" {
		flag.StringVar(&cfg.AccrualSystemAddress, "r", "", "адрес системы расчёта начислений")
	}
	flag.Parse()
	return &cfg, nil
}
