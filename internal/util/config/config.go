package config

import (
	"aulway/internal/util/model"
	"fmt"

	"github.com/vrischmann/envconfig"
)

func NewConfig() (model.Config, error) {
	var cfg model.Config
	if err := envconfig.Init(&cfg); err != nil {
		return model.Config{}, fmt.Errorf("get configs: %w", err)
	}

	return cfg, nil
}
