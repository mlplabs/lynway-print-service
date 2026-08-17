package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

func NewConfig() (*Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	return &cfg, nil
}
