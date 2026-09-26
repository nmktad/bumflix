package transcoder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ProcessVideoForHLSStreaming(ctx context.Context, inputPath string) (string, error) {
	outDir, err := os.MkdirTemp("", "bumflix-hls-*")
	if err != nil {
		return "", fmt.Errorf("creating output dir: %w", err)
	}

	playlist := filepath.Join(outDir, "index.m3u8")
	segmentPattern := filepath.Join(outDir, "segment_%03d.ts")

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", inputPath,
		"-c:v", "h264", "-c:a", "aac",
		"-hls_time", "6",
		"-hls_playlist_type", "vod",
		"-hls_segment_filename", segmentPattern,
		playlist,
		)

	cmd.Stderr = os.Stderr // ffmpeg logs progress to stderr; surface it while you're debugging

	if err := cmd.Run(); err != nil {
		os.RemoveAll(outDir)
		return "", fmt.Errorf("ffmpeg failed: %w", err)
	}

	return outDir, nil
}
