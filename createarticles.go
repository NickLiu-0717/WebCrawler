package main

import (
	"context"
	"time"

	database "github.com/NickLiu-0717/crawler/internal/database"
)

func (cfg *config) createArticles(url, title, content, catagory string, publishAt time.Time) error {
	_, err := cfg.db.CreateArticle(context.Background(), database.CreateArticleParams{
		Url:         url,
		Title:       title,
		Content:     content,
		Catagory:    catagory,
		PublishedAt: publishAt,
	})
	if err != nil {
		return err
	}
	return nil
}
