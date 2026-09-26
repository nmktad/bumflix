package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/nmktad/bumflix/internal/config"
	"github.com/nmktad/bumflix/internal/queue"
	"github.com/nmktad/bumflix/internal/storage"
	"github.com/nmktad/bumflix/internal/transcoder"
	"github.com/nmktad/bumflix/internal/utils"
)

// NOTE: just for simplicity there is s3Event from lambda library but felt weird to add it
type S3Event struct {
	Records []struct {
		S3 struct {
			Bucket struct {
				Name string `json:"name"`
			} `json:"bucket"`
			Object struct {
				Key string `json:"key"`
			} `json:"object"`
		} `json:"s3"`
	} `json:"Records"`
}

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

		var s3event S3Event

		if err := json.Unmarshal([]byte(*message.Body), &s3event); err != nil {
			fmt.Printf("Failed to unmarshal S3 body: %v", err)
			continue
		}

		var record = s3event.Records[0]

		filename := filepath.Base(record.S3.Object.Key)

		file, err := os.CreateTemp("", filename)

		if err != nil {
			fmt.Printf("Couldn't create file %v. Here's why: %v\n", filename, err)
			return 
		}

		defer os.Remove(file.Name())
		defer file.Close()

		storage.GetObject(context.Background(), record.S3.Object.Key, record.S3.Bucket.Name, file)

		absPath, err := filepath.Abs(file.Name())

		if err != nil {
			fmt.Printf("can't find the absolute path to file: %v", err)
			return
		}

		outDir, err := transcoder.ProcessVideoForHLSStreaming(context.Background(), absPath)
		if err != nil {
			fmt.Printf("hls video process failed: %v", err)
			return
		}

		storage.UploadDir(context.Background(), outDir, "bumflix-hls-segments", utils.MovieName(record.S3.Object.Key))

		queue.DeleteMessage(
			context.Background(), 
			fmt.Sprintf("http://sqs.%s.localhost.localstack.cloud:4566/000000000000/bumflix-transcode", cfg.AWSDefaultRegion), 
			*message,
			)
	}
}
