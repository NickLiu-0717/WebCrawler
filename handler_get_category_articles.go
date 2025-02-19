package main

import (
	"fmt"
	"net/http"
)

func (apicfg *apiConfig) handlerGetCategoryArticles(w http.ResponseWriter, r *http.Request) {
	cate := r.PathValue("category")
	fmt.Println(cate)

	dbArticles, err := apicfg.db.GetArticlesByCategory(r.Context(), cate)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't get articles by category from database", err)
		return
	}

	var articles []Article
	for _, dbArticle := range dbArticles {
		article := Article{
			ID:           dbArticle.ID,
			URL:          dbArticle.Url,
			Title:        dbArticle.Title,
			Content:      dbArticle.Content,
			Catagory:     dbArticle.Catagory,
			Created_at:   dbArticle.CreatedAt,
			Published_at: dbArticle.PublishedAt,
		}
		articles = append(articles, article)
	}
	respondWithJSON(w, http.StatusOK, articles)
}
