package main

import (
	"context"

	database "github.com/NickLiu-0717/crawler/internal/database"
)

func (cfg *config) createArticles(url, title, content string) error {
	_, err := cfg.db.CreateArticle(context.Background(), database.CreateArticleParams{
		Url:      url,
		Title:    title,
		Content:  content,
		Catagory: "news",
	})
	if err != nil {
		return err
	}
	return nil
}
