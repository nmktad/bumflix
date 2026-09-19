package storage

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/nmktad/bumflix/internal/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var client *s3.Client

func New() (*s3.Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("Couldn't read environment variables\n%v\n", err)
	}

	if client == nil {
		sdkcfg, err := awsconfig.LoadDefaultConfig(context.Background())
		if err != nil {
			return nil, fmt.Errorf("Couldn't load AWS default configuration\n%v\n", err)
		}

		client = s3.NewFromConfig(sdkcfg, func(o *s3.Options) {
			o.Region = cfg.AWSDefaultRegion
			o.BaseEndpoint = aws.String(cfg.AWSS3Endpoint)
		})
	}

	return client, nil
}
