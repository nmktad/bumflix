package benchmark

import (
	"encoding/csv"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	mu     sync.Mutex
	writer *csv.Writer
	file   *os.File
)

func Init(path string) error {
	if path == "" {
		return nil
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("opening bench log: %w", err)
	}

	stat, _ := f.Stat()
	file = f
	writer = csv.NewWriter(f)
	if stat.Size() == 0 {
		writer.Write([]string{"timestamp", "run_label", "movie", "download_ms", "transcode_ms", "upload_ms", "total_ms"})
		writer.Flush()
	}
	return nil
}

func Record(runLabel, movie string, downloadMs, transcodeMs, uploadMs, totalMs int64) {
	if writer == nil {
		return
	}

	mu.Lock()

	defer mu.Unlock()

	writer.Write([]string{
		time.Now().UTC().Format(time.RFC3339),
		runLabel, movie,
		fmt.Sprint(downloadMs), fmt.Sprint(transcodeMs), fmt.Sprint(uploadMs), fmt.Sprint(totalMs),
	})

	writer.Flush()
}
