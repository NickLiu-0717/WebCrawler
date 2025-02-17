package main

import "testing"

func TestCheckArticle(t *testing.T) {
	tests := []struct {
		name         string
		inputURL     string
		expected     bool
		errorContain string
	}{
		{
			name:     "bbc article",
			inputURL: "https://www.bbc.com/news/articles/cwy7e5ngwlzo",
			expected: true,
		},
		{
			name:     "bbc video",
			inputURL: "https://www.bbc.com/news/videos/c5y61gd0n3zo",
			expected: false,
		},
		{
			name:     "bbc catagory",
			inputURL: "https://www.bbc.com/news/",
			expected: false,
		},
		{
			name:     "nytime article",
			inputURL: "https://www.nytimes.com/2025/02/16/world/africa/in-latest-advance-rebels-in-congo-say-they-have-entered-a-key-city.html",
			expected: true,
		},
		{
			name:     "nytime crosswords",
			inputURL: "https://www.nytimes.com/crosswords",
			expected: false,
		},
		{
			name:     "pts.org",
			inputURL: "https://news.pts.org.tw/article/738105",
			expected: true,
		},
		{
			name:     "ebc.net",
			inputURL: "https://news.ebc.net.tw/news/society/471990",
			expected: true,
		},
		{
			name:     "ltn.com",
			inputURL: "https://news.ltn.com.tw/news/politics/breakingnews/4953899",
			expected: true,
		},
	}
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := checkArticle(tc.inputURL)
			if result != tc.expected {
				t.Errorf("Test %d %v FAIL - expect: %v, got: %v", i, tc.name, tc.expected, result)
			}
		})
	}
}
