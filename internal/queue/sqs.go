package queue

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/nmktad/bumflix/internal/config"
)

var client *sqs.Client

func New() (*sqs.Client, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("Couldn't read environment variables\n%v\n", err)
	}

	if client == nil {
		sdkcfg, err := awsconfig.LoadDefaultConfig(context.Background())
		if err != nil {
			return nil, fmt.Errorf("Couldn't load AWS default configuration\n%v\n", err)
		}

		client = sqs.NewFromConfig(sdkcfg, func(o *sqs.Options) {
			o.Region = cfg.AWSDefaultRegion
			o.BaseEndpoint = aws.String(cfg.AWSEndpoint)
		})
	}

	return client, nil
}

func ReadMessage(ctx context.Context, queueUrl string) (*types.Message, error) {
	c, err := New()
	if err != nil {
		return nil, fmt.Errorf("error encountered when creating a client: %w", err)
	}

	result, err := c.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueUrl),
		MaxNumberOfMessages: 1,
		WaitTimeSeconds: 20,
	})

	if err != nil {
		return nil, fmt.Errorf("error receiving messages: %w", err)
	}

	if len(result.Messages) == 0 {
		return nil, nil
	}

	return &result.Messages[0], nil
}

func DeleteMessage(ctx context.Context, queueUrl string, message types.Message) error {
	c, err := New()
	if err != nil {
		return fmt.Errorf("error encountered when creating a client: %w", err)
	}

	_, err = c.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),
		ReceiptHandle: message.ReceiptHandle,
	})
	if err != nil {
		return fmt.Errorf("couldn't delete message from queue %v: %w", queueUrl, err)
	}

	return nil
}
