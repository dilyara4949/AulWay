package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/vrischmann/envconfig"
)

func NewConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: No .env file found, using system environment variables")
	}

	var cfg Config
	if err := envconfig.Init(&cfg); err != nil {
		return Config{}, fmt.Errorf("get configs: %w", err)
	}

	return cfg, nil
}
