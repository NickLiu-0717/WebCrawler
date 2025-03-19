package main

import (
	"strings"
	"testing"

	"github.com/NickLiu-0717/crawler/crawl"
)

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		expected      string
		errorContains string
	}{
		{
			name:     "normalize 1",
			inputURL: "https://example.com/",
			expected: "example.com",
		},
		{
			name:     "normalize 2",
			inputURL: "http://google.com/",
			expected: "google.com",
		},
		{
			name:     "normalize 3",
			inputURL: "http://youtube.com/test1/",
			expected: "youtube.com/test1",
		},
		{
			name:          "Invalid URL",
			inputURL:      "://example.com",
			expected:      "",
			errorContains: "couldn't parse raw URL",
		},
	}
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			getResult, err := crawl.NormalizeURL(tc.inputURL)
			if err != nil && !strings.Contains(err.Error(), tc.errorContains) {
				t.Fatalf("Test %v - '%s' FAIL: error %v", i, tc.name, err)
			}
			if tc.expected != getResult {
				t.Errorf("Test %v - '%s' FAIL: expect: %s, got: %s", i, tc.name, tc.expected, getResult)
			}
		})
	}
}
