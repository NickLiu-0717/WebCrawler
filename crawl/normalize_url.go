package crawl

import (
	"fmt"
	"net/url"
	"strings"
)

func NormalizeURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("couldn't parse raw URL: %v", err)
	}

	var fullURL string
	fullURL = parsedURL.Host + parsedURL.Path
	fullURL = strings.ToLower(fullURL)
	fullURL = strings.TrimSuffix(fullURL, "/")

	return fullURL, nil
}
