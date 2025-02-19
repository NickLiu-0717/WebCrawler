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

	req, err := http.NewRequest("GET", parseURL.String(), nil)
	if err != nil {
		return "", fmt.Errorf("couldn't make new request: %v", err)
	}
	//Randomly pick User-Agent to prevent IP block
	req.Header.Set("User-Agent", getRandomUserAgent())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("couldn't get response from %s: %v", parseURL.String(), err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("error: URL %s response with code %v", parseURL.String(), resp.StatusCode)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		// return "", fmt.Errorf("error: response is not HTML, got: %v", contentType)
		return "", nil
	}
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("couldn't read response: %v", err)
	}

	return string(dat), nil
}
