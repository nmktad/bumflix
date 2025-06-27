package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/nmktad/bumflix/internal/s3client"
)

func main() {
	client, err := s3client.New()
	if err != nil {
		log.Fatalf("failed to create s3 client: %v", err)
	}

	hls_bucket := os.Getenv("AWS_S3_HLS_BUCKET_NAME")
	if hls_bucket == "" {
		log.Fatal("AWS_S3_HLS_BUCKET_NAME is not set")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/films", listVideosHandler(client, hls_bucket))
	r.Get("/film/{slug}", getMasterPlaylistHandler(client, hls_bucket))

	fmt.Println("Listening on :8080")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

func listVideosHandler(client *s3.Client, bucket string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket:    &bucket,
			Delimiter: aws.String("/"),
		})
		if err != nil {
			http.Error(w, "Could not list videos", http.StatusInternalServerError)
			return
		}

		var titles []string
		for _, prefix := range resp.CommonPrefixes {
			titles = append(titles, *prefix.Prefix)
		}

		w.Header().Set("Content-Type", "application/json")

		if _, err := w.Write([]byte(fmt.Sprintf(`{"videos": %q}`, titles))); err != nil {
			http.Error(w, "Failed to presign URL", http.StatusInternalServerError)
		}
	}
}

func getMasterPlaylistHandler(client *s3.Client, bucket string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slug := chi.URLParam(r, "slug")
		key := fmt.Sprintf("%s/master.m3u8", slug)

		presigned, err := s3.NewPresignClient(client).PresignGetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &key,
		}, s3.WithPresignExpires(1*time.Hour))
		if err != nil {
			http.Error(w, "Failed to presign URL", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, presigned.URL, http.StatusFound)
	}
}
