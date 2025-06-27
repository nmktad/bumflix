package transcoder

import (
	"os"
	"os/exec"
	"path/filepath"
)

type Variant struct {
	Name      string
	Width     int
	Height    int
	BitrateK  int
	OutputDir string
}

func GenerateHLSVariants(input, output string) error {
	args := []string{
		"-i", input,
		"-filter_complex",
		"[0:v]split=3[v1][v2][v3];" +
			"[v1]scale=w=1920:h=1080[v1out];" +
			"[v2]scale=w=1280:h=720[v2out];" +
			"[v3]scale=w=854:h=480[v3out]",
		"-map", "[v1out]", "-c:v:0", "libx264", "-b:v:0", "5000k", "-maxrate:v:0", "5350k", "-bufsize:v:0", "7500k",
		"-map", "[v2out]", "-c:v:1", "libx264", "-b:v:1", "2800k", "-maxrate:v:1", "2996k", "-bufsize:v:1", "4200k",
		"-map", "[v3out]", "-c:v:2", "libx264", "-b:v:2", "1400k", "-maxrate:v:2", "1498k", "-bufsize:v:2", "2100k",
		"-map", "a:0", "-c:a", "aac", "-b:a:0", "192k", "-ac", "2",
		"-map", "a:0", "-c:a", "aac", "-b:a:1", "128k", "-ac", "2",
		"-map", "a:0", "-c:a", "aac", "-b:a:2", "96k", "-ac", "2",
		"-f", "hls",
		"-hls_time", "10",
		"-hls_playlist_type", "vod",
		"-hls_flags", "independent_segments",
		"-hls_segment_type", "mpegts",
		"-hls_segment_filename", filepath.Join(output, "stream_%v/data%03d.ts"),
		"-master_pl_name", "master.m3u8",
		"-var_stream_map", "v:0,a:0 v:1,a:1 v:2,a:2",
		filepath.Join(output, "stream_%v/playlist.m3u8"),
	}

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
