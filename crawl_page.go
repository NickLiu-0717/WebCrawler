package main

import (
	"fmt"
	"net/url"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	parsedCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - couldn't parse current URL: %v\n", err)
		return
	}

	parseBaseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("Error - couldn't parse base URL: %v\n", err)
		return
	}

	if parseBaseURL.Hostname() != parsedCurrentURL.Hostname() {
		return
	}

	norCurrentURL, err := normalizeURL(parsedCurrentURL.String())
	if err != nil {
		fmt.Printf("Error - couldn't normalize current URL: %v\n", err)
		return
	}

	if count, ok := pages[norCurrentURL]; ok {
		pages[norCurrentURL] = count + 1
		return
	}
	pages[norCurrentURL] = 1

	fmt.Printf("Start crawling %s\n", rawCurrentURL)

	html, err := getHTML(parsedCurrentURL.String())
	if err != nil {
		fmt.Printf("Error - couldn't get HTML: %v\n", err)
		return
	}
	URLs, err := getURLsFromHTML(html, parsedCurrentURL.String())
	if err != nil {
		fmt.Printf("Error - couldn't get URL from HTML: %v\n", err)
		return
	}
	for _, nextURL := range URLs {
		crawlPage(rawBaseURL, nextURL, pages)
	}

}
