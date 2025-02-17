package main

import (
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
)

func extractArticles(htmlbody string) (string, string, error) {
	doc, err := htmlquery.Parse(strings.NewReader(htmlbody))
	if err != nil {
		return "", "", fmt.Errorf("couldn't parse html")
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

	if title == "" || len(content) < 100 {
		return "", "", fmt.Errorf("couldn't extrac article")
	}

	return title, content, nil

}
