package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestGetURLFromHTML(t *testing.T) {
	tests := []struct {
		name      string
		inputBody string
		inputURL  string
		expected  []string
	}{
		{name: "simple 1",
			inputBody: `
			<html lang="en">
			<head>
				<title>Links Example</title>
			</head>
			<body>
				<h2>Simple Links</h2>
				<a href="https://example.com">Visit Example</a>
				<br>
				<a href="https://golang.org">Go Language</a>
			</body>
			</html>
			`,
			inputURL: "https://example.com",
			expected: []string{
				"https://example.com",
				"https://golang.org",
			},
		},
		{name: "simple 2",
			inputBody: `
			<html>
				<body>
					<a href="/path/one">
						<span>Boot.dev</span>
					</a>
					<a href="https://other.com/path/one">
						<span>Boot.dev</span>
					</a>
				</body>
			</html>
			`,
			inputURL: "https://blog.boot.dev",
			expected: []string{
				"https://blog.boot.dev/path/one",
				"https://other.com/path/one",
			},
		},
		{name: "complex 1",
			inputBody: `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<title>Base URL Example</title>
			</head>
			<body>
				<h2>Navigation Links</h2>
				<a href="about">About Us</a>  <!-- Becomes https://example.com/about -->
				<br>
				<a href="contact">Contact</a>  <!-- Becomes https://example.com/contact -->
				<br>
				<a href="blog/post1">Blog Post 1</a>  <!-- Becomes https://example.com/blog/post1 -->

			</body>
			</html>
			`,
			inputURL: "https://example.com/",
			expected: []string{
				"https://example.com/about",
				"https://example.com/contact",
				"https://example.com/blog/post1",
			},
		},
	}

	cases := []struct {
		name          string
		inputURL      string
		inputBody     string
		expected      []string
		errorContains string
	}{
		{
			name:     "absolute URL",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="https://blog.boot.dev">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://blog.boot.dev"},
		},
		{
			name:     "relative URL",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://blog.boot.dev/path/one"},
		},
		{
			name:     "absolute and relative URLs",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href="/path/one">
			<span>Boot.dev</span>
		</a>
		<a href="https://other.com/path/one">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: []string{"https://blog.boot.dev/path/one", "https://other.com/path/one"},
		},
		{
			name:     "no href",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a>
			<span>Boot.dev></span>
		</a>
	</body>
</html>
`,
			expected: nil,
		},
		{
			name:     "bad HTML",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html body>
	<a href="path/one">
		<span>Boot.dev></span>
	</a>
</html body>
`,
			expected: []string{"https://blog.boot.dev/path/one"},
		},
		{
			name:     "invalid href URL",
			inputURL: "https://blog.boot.dev",
			inputBody: `
<html>
	<body>
		<a href=":\\invalidURL">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected: nil,
		},
		{
			name:     "handle invalid base URL",
			inputURL: `:\\invalidBaseURL`,
			inputBody: `
<html>
	<body>
		<a href="/path">
			<span>Boot.dev</span>
		</a>
	</body>
</html>
`,
			expected:      nil,
			errorContains: "couldn't parse base URL",
		},
	}
	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			getURLs, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil {
				t.Fatalf("Test %v - %v FAIL: unexpected error: %v", i, tc.name, err)
			}
			if !reflect.DeepEqual(getURLs, tc.expected) {
				t.Errorf("Test%v - %v FAIL: expect: %v, got: %v", i, tc.name, tc.expected, getURLs)
			}
		})
	}

	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			getURLs, err := getURLsFromHTML(tc.inputBody, tc.inputURL)
			if err != nil && !strings.Contains(err.Error(), tc.errorContains) {
				t.Errorf("Test %v - '%s' FAIL: unexpected error: %v", i, tc.name, err)
				return
			}
			if !reflect.DeepEqual(getURLs, tc.expected) {
				t.Errorf("Case %v - %v FAIL: expect: %v, got: %v", i, tc.name, tc.expected, getURLs)
			}
		})
	}
}
