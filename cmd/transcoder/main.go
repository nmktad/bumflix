package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nmktad/bumflix/internal/benchmark"
	"github.com/nmktad/bumflix/internal/config"
	databse "github.com/nmktad/bumflix/internal/database"
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

	// bench setup -- no-op if BENCH_LOG_PATH isn't set
	if err := benchmark.Init(os.Getenv("BENCH_LOG_PATH")); err != nil {
		fmt.Printf("warning: couldn't init bench log: %v\n", err)
	}
	runLabel := os.Getenv("BENCH_RUN_LABEL")

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, fmt.Sprintf(
		"postgresql://%s:%s@localhost:5432/%s?sslmode=disable",
		cfg.PGUser, cfg.PGPassword, cfg.PGDatabse,
		))
	if err != nil {
		panic(fmt.Errorf("can't connect to db: %v\n", err))
	}
	defer conn.Close(ctx)

	queries := databse.New(conn)
	queueUrl := fmt.Sprintf("http://sqs.%s.localhost.localstack.cloud:4566/000000000000/bumflix-transcode", cfg.AWSDefaultRegion)

	for {
		message, err := queue.ReadMessage(
			context.Background(), 
			queueUrl,
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

		if len(s3event.Records) == 0 {
			fmt.Println("event had no records, skipping")
			continue
		}

		var record = s3event.Records[0]

		rawKey := record.S3.Object.Key

		key, err := utils.DecodeS3Key(rawKey)
		if err != nil {
			fmt.Printf("failed to decode key %v: %v\n", rawKey, err)
			continue
		}

		filename := filepath.Base(key)

		file, err := os.CreateTemp("", filename)

		if err != nil {
			fmt.Printf("Couldn't create file %v. Here's why: %v\n", filename, err)
			return 
		}

		defer os.Remove(file.Name())
		defer file.Close()

		start := time.Now()

		// --- download ---
		t0 := time.Now()
		if err := storage.GetObject(context.Background(), key, record.S3.Bucket.Name, file); err != nil {
			fmt.Printf("failed to download %v: %v\n", key, err)
			file.Close()
			os.Remove(file.Name())
			continue
		}
		downloadMs := time.Since(t0).Milliseconds()

		absPath, err := filepath.Abs(file.Name())

		if err != nil {
			fmt.Printf("can't find the absolute path to file: %v", err)
			return
		}

		title := utils.MovieName(key)
		prefix := utils.Slugify(key)

		fmt.Println("transcode started")

		// --- transcode ---
		t0 = time.Now()
		outDir, err := transcoder.ProcessVideoForHLSStreaming(context.Background(), absPath)
		if err != nil {
			fmt.Printf("hls video process failed: %v", err)
			return
		}
		transcodeMs := time.Since(t0).Milliseconds()

		fmt.Println("transcode finished")

		t0 = time.Now()
		err = storage.UploadDir(context.Background(), outDir, "bumflix-hls-segments", prefix)

		uploadMs := time.Since(t0).Milliseconds()

		if err != nil {
			fmt.Printf("upload failed for %v: %v\n", prefix, err)
			continue
		}

		totalMs := time.Since(start).Milliseconds()
		benchmark.Record(runLabel, title, downloadMs, transcodeMs, uploadMs, totalMs)

		fmt.Println("creating the movie entry in db")
		movie, err := queries.CreateMovie(ctx, databse.CreateMovieParams{
			Title: title,
			Prefix: prefix,
			Description: pgtype.Text{
				String: "an amazing movie",
				Valid: true,
			},
		})

		if err != nil {
			fmt.Printf("error creating a movie %v", err)
			return 
		} else {
			fmt.Println("created the movie entry in db")
			fmt.Printf("movie: %v", movie)
		}

		queue.DeleteMessage(
			context.Background(), 
			fmt.Sprintf("http://sqs.%s.localhost.localstack.cloud:4566/000000000000/bumflix-transcode", cfg.AWSDefaultRegion), 
			*message,
			)

	 
		fmt.Println("finished the process movie ready")
	}
}
