package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/nmktad/bumflix/internal/config"
	"github.com/nmktad/bumflix/internal/storage"
)

func main() {
	sourceFile := flag.String("file", "", "path to a sample video")
	count := flag.Int("n", 20, "number of jobs to queue")
	flag.Parse()

	if *sourceFile == "" {
		fmt.Println("usage: benchload -file testdata/sample.mp4 -n 20")
		os.Exit(1)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	c, err := storage.New()
	if err != nil {
		panic(err)
	}

	start := time.Now()
	for i := 0; i < *count; i++ {
		f, err := os.Open(*sourceFile)
		if err != nil {
			panic(err)
		}

		key := fmt.Sprintf("movies/bench-%04d.mp4", i)
		_, err = c.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(cfg.SourceBucketName),
			Key:    aws.String(key),
			Body:   f,
		})
		f.Close()

		if err != nil {
			fmt.Printf("failed to upload %v: %v\n", key, err)
			continue
		}
		fmt.Printf("queued %v\n", key)
	}
	fmt.Printf("queued %d jobs in %v\n", *count, time.Since(start))
}
