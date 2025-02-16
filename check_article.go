package main

import (
	"regexp"
	"strings"
)

func checkArticle(url string) bool {
	// 常見的文章 URL 格式（可依照實際網站調整）
	articlePatterns := []string{
		`/\d{4}/\d{2}/\d{2}/`, // YYYY/MM/DD 格式 (e.g., example.com/2024/02/16/title/)
		`/posts/\d+`,          // post ID 格式 (e.g., example.com/posts/12345)
		`/news/\d+`,           // e.g., example.com/news/67890
		`/blog/[\w-]+`,        // 部落格標題 (e.g., example.com/blog/my-article-title)
		`/news/articles/[a-zA-Z0-9]+$`,
	}

	// 不要的 URL（分類、標籤、登入等）
	nonArticleKeywords := []string{"category", "tag", "archive", "login", "feed", "rss"}

	// 檢查是否符合文章格式
	for _, pattern := range articlePatterns {
		match, _ := regexp.MatchString(pattern, url)
		if match {
			return true
		}
	}

	// 排除非文章類的 URL
	for _, keyword := range nonArticleKeywords {
		if strings.Contains(url, keyword) {
			return false
		}
	}

	return false
}
