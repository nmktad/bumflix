package utils

import (
	"fmt"
	"path"
	"strings"

	"github.com/nmktad/bumflix/internal/config"
	databse "github.com/nmktad/bumflix/internal/database"
)

func MovieName(key string) string {
	base := path.Base(key)
	return strings.TrimSuffix(base, path.Ext(base))
}

func PlaybackURL(movie databse.Movie) string {
	cfg, err := config.LoadConfig()
	if err != nil {
		return ""
	}

	return fmt.Sprintf("http://localhost:4566/%s/%s/index.m3u8", cfg.DestBucketName, movie.Prefix)
}
