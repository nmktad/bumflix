package utils

import (
	"path"
	"strings"
)

func MovieName(key string) string {
	base := path.Base(key)
	return strings.TrimSuffix(base, path.Ext(base))
}

