package utils

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/nmktad/bumflix/internal/config"
	databse "github.com/nmktad/bumflix/internal/database"

	"net/url"
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

var nonSlugChars = regexp.MustCompile(`[^a-z0-9]+`)

func Slugify(s string) string {
	s = strings.ToLower(s)
	s = nonSlugChars.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

func DecodeS3Key(key string) (string, error) {
	return url.QueryUnescape(strings.ReplaceAll(key, "+", "%20"))
}
