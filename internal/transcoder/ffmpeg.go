package transcoder

import (
	"fmt"

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
