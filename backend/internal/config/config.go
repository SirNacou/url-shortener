package config

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port      string `env:"PORT" envDefault:"8080"`
	TableName string `env:"TABLE_NAME,required"`
	AppEnv    string `env:"APP_ENV" envDefault:"production"`

	AWSConfig aws.Config `json:"-"`
}

func Load(ctx context.Context) (*Config, error) {
	var cfg Config

	// 1. Parse plain environment variables
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse env: %w", err)
	}

	// 2. Load AWS SDK configuration explicitly
	awsCfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	cfg.AWSConfig = awsCfg
	return &cfg, nil
}
