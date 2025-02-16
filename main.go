package main

import (
	"fmt"
	"net/url"
	"os"
	"sort"
)

type urlCount struct {
	url   string
	count int
}

func main() {
	cmdArg := os.Args

	if len(cmdArg) == 1 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(cmdArg) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	} else if len(cmdArg) == 2 {
		_, err := url.Parse(cmdArg[1])
		if err != nil {
			fmt.Printf("couldn't parse url %v : %v\n", cmdArg[1], err)
			os.Exit(1)
		}
	}

	rawURL := cmdArg[1]
	pages := make(map[string]int)
	crawlPage(rawURL, rawURL, pages)

	// for url, count := range pages {
	// 	fmt.Printf("%s: %d\n", url, count)
	// }

	var urlCounts []urlCount
	for url, count := range pages {
		urlCounts = append(urlCounts, urlCount{url, count})
	}

	// Sort the slice (you can use sort.Slice)
	sort.Slice(urlCounts, func(i, j int) bool {
		return urlCounts[i].count > urlCounts[j].count // '>' for descending order
	})

	// Print sorted results
	for _, uc := range urlCounts {
		fmt.Printf("%s: %d\n", uc.url, uc.count)
	}
}
