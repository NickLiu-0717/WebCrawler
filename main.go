package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const defaultMaxConcurrency = 5
const defaultMaxPages = 10
const defaultMaxDepth = 5

// var newsPages []string = []string{
// 	"https://www.bbc.com/",
// 	"https://edition.cnn.com/",
// 	"https://www.nytimes.com/",
// 	"https://www.forbes.com/",
// 	"https://news.ebc.net.tw",
// 	"https://news.ltn.com.tw",
// 	"https://news.pts.org.tw/",
// }

func main() {
	var maxConcurrency int
	var maxPages int
	var maxDepth int
	var err error
	cmdArg := os.Args

	// if len(cmdArg) == 1 {
	// 	fmt.Println("Error - need arguments: website maxConcurrency maxPages")
	// 	os.Exit(1)
	// }

	// if len(cmdArg) > 4 {
	// 	fmt.Println("Error - too many arguments provided")
	// 	os.Exit(1)
	// }

	// if len(cmdArg) == 2 {
	// 	maxConcurrency = defaultMaxConcurrency
	// 	maxPages = defaultMaxPages
	// }

	// if len(cmdArg) == 3 {
	// 	maxConcurrency, err = strconv.Atoi(cmdArg[2])
	// 	if err != nil {
	// 		fmt.Println("Error - couldn't convert input string to integer")
	// 		return
	// 	}
	// 	maxPages = defaultMaxPages
	// }

	// if len(cmdArg) == 4 {
	// 	maxConcurrency, err = strconv.Atoi(cmdArg[2])
	// 	if err != nil {
	// 		fmt.Println("Error - couldn't convert input string to integer")
	// 		return
	// 	}
	// 	maxPages, err = strconv.Atoi(cmdArg[3])
	// 	if err != nil {
	// 		fmt.Println("Error - couldn't convert input string to integer")
	// 		return
	// 	}
	// }

	switch len(cmdArg) {
	case 1:
		fmt.Println("Error - need arguments: website maxConcurrency maxPages")
		os.Exit(1)
	case 2:
		maxConcurrency = defaultMaxConcurrency
		maxPages = defaultMaxPages
		maxDepth = defaultMaxDepth
	case 3:
		maxConcurrency, err = strconv.Atoi(cmdArg[2])
		if err != nil {
			fmt.Println("Error - couldn't convert input string to integer")
			return
		}
		maxPages = defaultMaxPages
		maxDepth = defaultMaxDepth
	case 4:
		maxConcurrency, err = strconv.Atoi(cmdArg[2])
		if err != nil {
			fmt.Println("Error - couldn't convert input string to integer")
			return
		}
		maxPages, err = strconv.Atoi(cmdArg[3])
		if err != nil {
			fmt.Println("Error - couldn't convert input string to integer")
			return
		}
		maxDepth = defaultMaxDepth
	case 5:
		maxConcurrency, err = strconv.Atoi(cmdArg[2])
		if err != nil {
			fmt.Println("Error - couldn't convert input string to integer")
			return
		}
		maxPages, err = strconv.Atoi(cmdArg[3])
		if err != nil {
			fmt.Println("Error - couldn't convert input string to integer")
			return
		}
		maxDepth, err = strconv.Atoi(cmdArg[4])
		if err != nil {
			fmt.Println("Error - couldn't convert input string to integer")
			return
		}
	default:
		fmt.Println("Error - too many arguments provided")
		os.Exit(1)
	}

	cfg, err := configure(cmdArg[1], maxConcurrency, maxPages, maxDepth)
	if err != nil {
		fmt.Printf("Error - couldn't configure: %v\n", err)
		return
	}

	group, err := checkRobotsTxt(cfg.baseURL.String())
	if err != nil {
		fmt.Printf("Error - couldn't check robots.txt: %v\n", err)
		return
	}
	cfg.robotGroup = group

	cfg.wg.Add(1)
	go cfg.crawlPage(cfg.baseURL.String(), 1)
	cfg.wg.Wait()

	printPages(cfg.pages, strings.TrimSuffix(cfg.baseURL.String(), "/"))
	for _, article := range cfg.articles {
		fmt.Println(article)
	}
}
