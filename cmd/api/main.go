package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	s3Client "github.com/nmktad/bumflix/pkg/s3"
)

func main() {
	http.HandleFunc("/", serveHLS)

	// Route to serve video content
	http.HandleFunc("/movies/", serveHLS)

	fmt.Printf("Starting server on port %d...\n", 8080)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", 8080), nil))
}

// serveHLS fetches HLS files from MinIO and serves them with correct headers
func serveHLS(w http.ResponseWriter, r *http.Request) {
	u, err := s3Client.New(os.Getenv("AWS_S3_BUCKET_NAME"))
	if err != nil {
		log.Fatal("couldn't connect to s3")
	}

	w.Header().Set("Access-Control-Allow-Origin", string("http://localhost:3000"))
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// If it's a preflight OPTIONS request, respond with a 200
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extract file path from the request
	objectName := strings.TrimPrefix(r.URL.Path, "/movies/")

	// Generate a presigned URL for secure access
	presignedURL, err := u.GetPresignedURL(objectName)
	if err != nil {
		http.Error(w, "Error generating presigned URL", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, presignedURL, http.StatusFound)
}
