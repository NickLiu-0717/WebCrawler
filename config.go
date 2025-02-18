package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"sync"
	"time"

	database "github.com/NickLiu-0717/crawler/internal/database"
	"github.com/temoto/robotstxt"
)

// type urlCount struct {
// 	url   string
// 	count int
// }

type Article struct {
	title   string
	content string
}

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
	robotGroup         *robotstxt.Group
	maxDepth           int
	articles           map[string]Article
	db                 *database.Queries
}

func configure(rawBaseURL string, maxConcurrency, maxPages, maxDepth int) (*config, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse url %v : %v", rawBaseURL, err)
	}
	return &config{
		pages:              make(map[string]int),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
		maxDepth:           maxDepth,
		articles:           make(map[string]Article),
	}, nil
}

// func printPages(pages map[string]int, baseURL string) {
// 	fmt.Println("=============================")
// 	fmt.Println("REPORT for", baseURL)
// 	fmt.Println("=============================")
// 	var urlCounts []urlCount
// 	for url, count := range pages {
// 		urlCounts = append(urlCounts, urlCount{url, count})
// 	}

// 	// Sort the slice (you can use sort.Slice)
// 	sort.Slice(urlCounts, func(i, j int) bool {
// 		return urlCounts[i].count > urlCounts[j].count // '>' for descending order
// 	})

// 	// Print sorted results
// 	for _, uc := range urlCounts {
// 		if uc.count > 1 {
// 			fmt.Printf("Found %d internal links to %s\n", uc.count, uc.url)
// 		} else {
// 			fmt.Printf("Found %d internal link to %s\n", uc.count, uc.url)
// 		}
// 		// fmt.Printf("%s: %d\n", uc.url, uc.count)
// 	}
// }

func randomSleep(min, max int) {
	delay := rand.Intn(max-min) + min
	time.Sleep(time.Duration(delay) * time.Millisecond)
}
