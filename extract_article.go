package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/araddon/dateparse"
	"golang.org/x/net/html"
)

func extractArticles(htmlbody string) (string, string, time.Time, error) {
	doc, err := htmlquery.Parse(strings.NewReader(htmlbody))
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("couldn't parse html")
	}

	//Find title from tag of meta[@property="og:title"] or h1
	var title string
	metatitle := htmlquery.FindOne(doc, `//meta[@property="og:title"]/@content`)
	if metatitle != nil {
		title = htmlquery.SelectAttr(metatitle, "content")
	} else {
		h1Title := htmlquery.FindOne(doc, "//h1")
		if h1Title != nil {
			title = htmlquery.InnerText(h1Title)
		}
	}
	title = strings.TrimSpace(title)

	//Find paragraphs from tag of p
	var paragraphs []string

	articleNode := htmlquery.FindOne(doc, "//article")
	if articleNode != nil {
		pNodes := htmlquery.Find(articleNode, "//p")
		for _, pNode := range pNodes {
			text := strings.TrimSpace(htmlquery.InnerText(pNode))
			paragraphs = append(paragraphs, text)
		}
	} else {
		pNodes := htmlquery.Find(doc, "//p")
		for _, pNode := range pNodes {
			text := strings.TrimSpace(htmlquery.InnerText(pNode))
			paragraphs = append(paragraphs, text)
		}
	}

	content := strings.Join(paragraphs, "\n\n")
	//extract image url (future feature)
	// imageNode := htmlquery.FindOne(doc, `//meta[@property="og:image"]/@content`)
	// imageURL := ""
	// if imageNode != nil {
	// 	imageURL = htmlquery.SelectAttr(imageNode, "content")
	// }
	publishedAt := extractPublishedTime(doc)
	parsedPub, err := dateparse.ParseAny(publishedAt)
	if err != nil {
		return "", "", time.Time{}, fmt.Errorf("couldn't parse published time")
	}

	if title == "" || len(content) < 100 {
		return "", "", time.Time{}, fmt.Errorf("couldn't extrac article")
	}

	return title, content, parsedPub, nil

}

func extractPublishedTime(doc *html.Node) string {
	xpathQueries := []string{
		"//meta[@property='article:published_time']/@content", // OpenGraph metadata
		"//meta[@name='datePublished']/@content",              // Schema.org metadata
		"//meta[@name='publish_date']/@content",               // Custom metadata
		"//meta[@name='date']/@content",                       // Generic date metadata
		"//time/@datetime",                                    // <time> tag with datetime attribute
		"//time/text()",                                       // <time> tag text content
		"//*[contains(@class, 'date')]/text()",                // Generic class containing "date"
		"//*[contains(@class, 'published')]/text()",           // Class containing "published"
		"//*[contains(@class, 'post-date')]/text()",           // Class containing "post-date"
		"//*[contains(@class, 'timestamp')]/text()",           // Class containing "timestamp"
	}
	for _, query := range xpathQueries {
		node := htmlquery.FindOne(doc, query)
		if node != nil {
			return htmlquery.InnerText(node)
		}
	}
	return "Published date not found"
}

// func tryParseTime(timeStr string) (time.Time, error) {
// 	// List of common time formats
// 	layouts := []string{
// 		"2006-01-02T15:04:05-07:00",
// 		"2006-01-02T15:04:05.999Z",
// 	}

// 	var parsedTime time.Time
// 	var err error
// 	for _, layout := range layouts {
// 		parsedTime, err = time.Parse(layout, timeStr)
// 		if err == nil {
// 			return parsedTime, nil // Successfully parsed
// 		}
// 	}

// 	return time.Time{}, fmt.Errorf("could not parse time: %s", timeStr)
// }
