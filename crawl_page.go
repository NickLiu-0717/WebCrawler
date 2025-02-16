package main

import (
	"fmt"
	"net/url"
)

func (cfg *config) crawlPage(rawCurrentURL string, depth int) {

	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	if depth > cfg.maxDepth {
		return
	}

	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}
	cfg.mu.Unlock()

	parsedCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - couldn't parse current URL: %v\n", err)
		return
	}

	if cfg.baseURL.Hostname() != parsedCurrentURL.Hostname() {
		return
	}

	norCurrentURL, err := normalizeURL(parsedCurrentURL.String())
	if err != nil {
		fmt.Printf("Error - couldn't normalize current URL: %v\n", err)
		return
	}

	if cfg.robotGroup != nil {
		if !cfg.robotGroup.Test(parsedCurrentURL.Path) {
			fmt.Printf("URL %s is not allowed to crawl\n", parsedCurrentURL.String())
			return
		}
	}

	if checkArticle(parsedCurrentURL.String()) {
		cfg.articles = append(cfg.articles, parsedCurrentURL.String())
	}

	cfg.mu.Lock()
	if _, ok := cfg.pages[norCurrentURL]; ok {
		cfg.pages[norCurrentURL]++
		cfg.mu.Unlock()
		return
	}
	cfg.pages[norCurrentURL] = 1
	cfg.mu.Unlock()

	fmt.Printf("Start crawling %s\n", rawCurrentURL)

	randomSleep(1000, 3000)
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
		cfg.wg.Add(1)
		go cfg.crawlPage(nextURL, depth+1)
	}

}
