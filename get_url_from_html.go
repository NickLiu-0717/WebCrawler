package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse base URL: %v", err)
	}

	// doc, err := html.Parse(strings.NewReader(htmlBody))
	// if err != nil {
	// 	return nil, fmt.Errorf("couldn't parse html: %v", err)
	// }
	// var urls []string
	// for n := range doc.Descendants() {
	// 	if n.Type == html.ElementNode && n.DataAtom == atom.A {			//old version, use doc.Descendants
	// 		for _, anchor := range n.Attr {
	// 			if anchor.Key == "href" {
	// 				href, err := url.Parse(anchor.Val)
	// 				if err != nil {
	// 					fmt.Printf("couldn't parse href %v: %v\n", anchor.Val, err)
	// 					continue
	// 				}
	// 				revolveURL := baseURL.ResolveReference(href)
	// 				urls = append(urls, revolveURL.String())
	// 				break
	// 			}
	// 		}
	// 	}
	// }

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
