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

	baseOutput := filepath.Join("output", movieID)

	if _, err := os.Stat(baseOutput); os.IsNotExist(err) {
		err := os.MkdirAll(baseOutput, 0755)
		if err != nil {
			return fmt.Errorf("failed to create output dir %s: %w", "output", err)
		}
	}

	variants := []transcoder.Variant{
		{Name: "240p", Width: 426, Height: 240, BitrateK: 400, OutputDir: filepath.Join(baseOutput, "240p")},
		{Name: "360p", Width: 640, Height: 360, BitrateK: 600, OutputDir: filepath.Join(baseOutput, "360p")},
		{Name: "480p", Width: 854, Height: 480, BitrateK: 800, OutputDir: filepath.Join(baseOutput, "480p")},
		{Name: "720p", Width: 1280, Height: 720, BitrateK: 1500, OutputDir: filepath.Join(baseOutput, "720p")},
		{Name: "1080p", Width: 1920, Height: 1080, BitrateK: 3000, OutputDir: filepath.Join(baseOutput, "1080p")},
		// {Name: "2160p", Width: 3840, Height: 2160, BitrateK: 8000, OutputDir: filepath.Join(baseOutput, "2160p")},
	}

	err := transcoder.GenerateHLSVariants(inputPath, variants)
	if err != nil {
		log.Fatal("Transcoding failed:", err)
	}

	for _, v := range variants {
		if _, err := os.Stat(fmt.Sprintf("output/la-jetee/%s", v.Name)); os.IsNotExist(err) {
			err := os.MkdirAll(fmt.Sprintf("output/la-jetee/%s", v.Name), 0755)
			if err != nil {
				return fmt.Errorf("failed to create output dir %s: %w", "output", err)
			}
		}

		err := s3client.UploadDir(v.OutputDir, fmt.Sprintf("la-jetee/%s", v.Name), os.Getenv("AWS_S3_HLS_BUCKET_NAME"))
		if err != nil {
			log.Fatal("Upload failed:", err)
		}
	}

	fmt.Println("movie has been processed perfectly")

	return nil
}
