package crawl

import (
	"context"
	"time"

	database "github.com/NickLiu-0717/crawler/internal/database"
)

func (cfg *CrawlConfig) createArticles(url, title, content, catagory string, publishAt time.Time) error {
	_, err := cfg.Config.Db.CreateArticle(context.Background(), database.CreateArticleParams{
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
