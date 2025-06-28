package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/nmktad/bumflix/internal/s3client"
	"github.com/nmktad/bumflix/src/film"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	// err := ingest.IngestExample()
	// if err != nil {
	// 	log.Fatalf("failed to create s3 client: %v", err)
	// }

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
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // frontend origin
		AllowedMethods:   []string{"GET", "POST", "HEAD", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Location"}, // optional
		AllowCredentials: true,
		MaxAge:           300, // 5 min
	}))

	r.Get("/films", listVideosHandler(client, hls_bucket))
	r.Head("/film/{slug}", getMasterPlaylistHandler(client, hls_bucket))
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

		fmt.Println(resp)
		if err != nil {
			http.Error(w, "Could not list videos", http.StatusNotFound)
			return
		}

		caser := cases.Title(language.English)

		films := make([]film.Film, 0)

		for _, prefix := range resp.CommonPrefixes {
			if prefix.Prefix == nil {
				continue
			}
			slug := strings.TrimSuffix(*prefix.Prefix, "/")
			title := strings.ReplaceAll(slug, "-", " ") // naive title formatting

			films = append(films, film.Film{
				Title: caser.String(title),
				Slug:  slug,
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string][]film.Film{"videos": films}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
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
