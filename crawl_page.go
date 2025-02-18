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
	//Normalize the path
	norCurrentURL, err := normalizeURL(parsedCurrentURL.String())
	if err != nil {
		fmt.Printf("Error - couldn't normalize current URL: %v\n", err)
		return
	}
	//Test if the path has robots.txt restriction
	if cfg.robotGroup != nil {
		if !cfg.robotGroup.Test(parsedCurrentURL.Path) {
			fmt.Printf("URL %s is not allowed to crawl\n", parsedCurrentURL.String())
			return
		}
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
	//Random Sleep to simulate human behavior
	randomSleep(1000, 2000)
	html, err := getHTML(parsedCurrentURL.String())
	if err != nil {
		fmt.Printf("Error - couldn't get HTML: %v\n", err)
		return
	}
	//Check if the path is article or not, if article, extract title and content then return, no more crawling
	if checkArticle(parsedCurrentURL.String()) {
		gottitle, gotcontent, gotpubAt, err := extractArticles(html)
		if err != nil {
			fmt.Printf("Error - couldn't extract title and content from HTML: %v\n", err)
			return
		}
		// cfg.articles[norCurrentURL] = Article{
		// 	title:   gottitle,
		// 	content: gotcontent,
		// }
		err = cfg.createArticles(norCurrentURL, gottitle, gotcontent, gotpubAt.UTC())
		if err != nil {
			fmt.Printf("couldn't create article to database: %v\n", err)
			return
		}
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
