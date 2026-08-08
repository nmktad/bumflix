package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type Config struct {
	SourceBucketName string `env:"SOURCE_BUCKET_NAME"`
	DestBucketName string `env:"DEST_BUCKET_NAME"`
	AWSDefaultRegion string `env:"AWS_DEFAULT_REGION"`
	AWSS3Endpoint string `env:"AWS_S3_ENDPOINT"`
	AWSEndpoint string `env:"AWS_ENDPOINT"`
}

var cfg *Config

func LoadConfig() (*Config, error) {
	if cfg == nil {
		cfg = &Config{}

		if err := env.Parse(cfg); err != nil {
			return nil, fmt.Errorf("parsing environment variables: %w", err)
		}
	}

	return cfg, nil
}
