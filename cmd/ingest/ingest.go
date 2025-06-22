package ingest

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nmktad/bumflix/pkg/ffmpeg"
	S3 "github.com/nmktad/bumflix/pkg/s3"
)

func IngestExample(movieTitle, path string) error {
	movieID := strings.ToLower(
		strings.ReplaceAll(
			strings.TrimSuffix(movieTitle, filepath.Ext(path)),
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

	variants := []ffmpeg.Variant{
		{Name: "240p", Width: 426, Height: 240, BitrateK: 400, OutputDir: filepath.Join(baseOutput, "240p")},
		{Name: "360p", Width: 640, Height: 360, BitrateK: 600, OutputDir: filepath.Join(baseOutput, "360p")},
		{Name: "480p", Width: 854, Height: 480, BitrateK: 800, OutputDir: filepath.Join(baseOutput, "480p")},
		{Name: "720p", Width: 1280, Height: 720, BitrateK: 1500, OutputDir: filepath.Join(baseOutput, "720p")},
		{Name: "1080p", Width: 1920, Height: 1080, BitrateK: 3000, OutputDir: filepath.Join(baseOutput, "1080p")},
		// {Name: "2160p", Width: 3840, Height: 2160, BitrateK: 8000, OutputDir: filepath.Join(baseOutput, "2160p")},
	}

	err := ffmpeg.GenerateHLSVariants(path, variants)
	if err != nil {
		log.Fatal("Transcoding failed:", err)
	}

	err = ffmpeg.WriteMasterPlaylist(baseOutput, variants)
	if err != nil {
		log.Fatal("Failed to write master.m3u8:", err)
	}

	u, err := S3.New(os.Getenv("AWS_S3_BUCKET_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range variants {
		if _, err := os.Stat(fmt.Sprintf("output/%s/%s", movieID, v.Name)); os.IsNotExist(err) {
			err := os.MkdirAll(fmt.Sprintf("output/%s/%s", movieID, v.Name), 0755)
			if err != nil {
				return fmt.Errorf("failed to create output dir %s: %w", "output", err)
			}
		}

		err := u.UploadDir(v.OutputDir, fmt.Sprintf("movies/%s/%s", movieID, v.Name))
		if err != nil {
			log.Fatal("Upload failed:", err)
		}

		localPath := filepath.Join(baseOutput, "master.m3u8")
		remoteKey := filepath.Join(movieID, "master.m3u8") // e.g., movies/it-happened-one-night_1934/master.m3u8

		err = u.UploadFile(localPath, remoteKey, "application/vnd.apple.mpegurl")
		if err != nil {
			log.Fatalf("Failed to upload master.m3u8: %v", err)
		}

	}

	fmt.Println("movie has been processed perfectly")

	return nil
}
