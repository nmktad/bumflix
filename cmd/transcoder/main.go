package main

import (
	"context"
	"fmt"

	"github.com/nmktad/bumflix/internal/config"
	"github.com/nmktad/bumflix/internal/queue"
)

func main() {
	cfg, err := config.LoadConfig() 
	if err != nil {
		panic(fmt.Errorf("Couldn't read environment variables\n%v\n", err))
	}

	for {
		message, err := queue.ReadMessage(
			context.Background(), 
			fmt.Sprintf("http://sqs.%s.localhost.localstack.cloud:4566/000000000000/bumflix-transcode", cfg.AWSDefaultRegion),
		)

		if err != nil {
			panic(fmt.Errorf("error while reading message from queue\n%v\n", err))
		}

		if message == nil {
			continue
		}

		fmt.Printf("read: %v", *message.Body)

		queue.DeleteMessage(
			context.Background(), 
			fmt.Sprintf("http://sqs.%s.localhost.localstack.cloud:4566/000000000000/bumflix-transcode", cfg.AWSDefaultRegion), 
			*message,
			)
	}
}
