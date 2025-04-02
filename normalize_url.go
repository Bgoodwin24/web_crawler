package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(rawUrl string) (string, error) {
	parsed, err := url.Parse(rawUrl)
	if err != nil {
		return "", fmt.Errorf("error parsing url: %v", err)
	}

	normalized := parsed.Host + parsed.Path

	trimmed := strings.TrimSuffix(normalized, "/")
	return trimmed, nil
}
