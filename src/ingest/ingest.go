package ingest

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nmktad/bumflix/internal/s3client"
	"github.com/nmktad/bumflix/internal/transcoder"
)

func IngestExample() error {
	movieFile := "Hotel-Chevalier_2007.mp4"
	inputPath := filepath.Join("local-movies", movieFile)

	movieID := strings.ToLower(
		strings.ReplaceAll(
			strings.TrimSuffix(movieFile, filepath.Ext(movieFile)),
			" ", "-",
		),
	)

	// Output structure: /tmp/processed/<movieID>
	baseOutput := filepath.Join(os.TempDir(), "processed", movieID)

	// Ensure output directory exists
	if err := os.MkdirAll(baseOutput, 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	log.Printf("Starting transcode for %s → %s", inputPath, baseOutput)

	// Run transcoding
	if err := transcoder.GenerateHLSVariants(inputPath, baseOutput); err != nil {
		return fmt.Errorf("transcode failed: %w", err)
	}

	log.Printf("Transcoding complete for %s", movieID)

	if err := s3client.UploadDir(baseOutput, movieID, os.Getenv("AWS_S3_HLS_BUCKET_NAME")); err != nil {
		return fmt.Errorf("failed to upload HLS output: %w", err)
	}

	log.Printf("Upload complete for %s", movieID)

	return nil
}
