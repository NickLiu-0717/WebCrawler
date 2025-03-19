package crawl

import (
	"context"
	"time"

	"github.com/NickLiu-0717/crawler/config"
	database "github.com/NickLiu-0717/crawler/internal/database"
)

type CrawlConfig struct {
	Config *config.Config
}

func (cfg *CrawlConfig) CreateArticles(url, title, content, catagory string, publishAt time.Time) error {
	_, err := cfg.Config.Db.CreateArticle(context.Background(), database.CreateArticleParams{
		Url:         url,
		Title:       title,
		Content:     content,
		Catagory:    catagory,
		PublishedAt: publishAt,
	})
	return err
}
