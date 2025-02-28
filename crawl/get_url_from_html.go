package crawl

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
)

func GetURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	doc, err := htmlquery.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, fmt.Errorf("couldn't parse html body: %v", err)
	}

	var urls []string
	nodes := htmlquery.Find(doc, "//a[@href]")
	for _, node := range nodes {
		href := htmlquery.SelectAttr(node, "href")
		parsedHref, err := url.Parse(href)
		if err != nil {
			fmt.Printf("couldn't parse href: %v\n", err)
			continue
		}
		fullURL := baseURL.ResolveReference(parsedHref)
		urls = append(urls, fullURL.String())
	}

	return urls, nil
}
