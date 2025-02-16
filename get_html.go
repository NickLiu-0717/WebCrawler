package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	parseURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("couldn't parse raw URL: %v", err)
	}
	resp, err := http.Get(parseURL.String())
	if err != nil {
		return "", fmt.Errorf("couldn't get response from %s: %v", parseURL.String(), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("error: response with code %v", resp.StatusCode)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("error: response is not HTML, got: %v", contentType)
	}
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("couldn't read response: %v", err)
	}

	return string(dat), nil
}
