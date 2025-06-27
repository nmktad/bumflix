package transcoder

import (
	"fmt"
	"os"
	"path/filepath"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type Variant struct {
	Name      string
	Width     int
	Height    int
	BitrateK  int
	OutputDir string
}

func GenerateHLSVariants(input string, variants []Variant) error {
	for _, v := range variants {
		outputPath := fmt.Sprintf("%s/index.m3u8", v.OutputDir)

		err := ffmpeg.Input(input).
			Filter("scale", ffmpeg.Args{fmt.Sprint(v.Width), fmt.Sprint(v.Height)}).
			Output(outputPath,
				ffmpeg.KwArgs{
					"c:v":                  "libx264",
					"b:v":                  fmt.Sprintf("%dk", v.BitrateK),
					"c:a":                  "aac",
					"hls_time":             6,
					"hls_playlist_type":    "vod",
					"hls_segment_filename": fmt.Sprintf("%s/segment_%%03d.ts", v.OutputDir),
					"f":                    "hls",
				}).
			OverWriteOutput().
			Run()
		if err != nil {
			return fmt.Errorf("failed for %s: %w", v.Name, err)
		}
	}
	return nil
}

func WriteMasterPlaylist(outputDir string, variants []Variant) error {
	var builder string
	builder += "#EXTM3U\n"

	for _, v := range variants {
		line := fmt.Sprintf(
			"#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d\n%s/index.m3u8\n",
			v.BitrateK*1000, v.Width, v.Height, v.Name,
		)
		builder += line
	}

	masterPath := filepath.Join(outputDir, "master.m3u8")
	return os.WriteFile(masterPath, []byte(builder), 0644)
}
