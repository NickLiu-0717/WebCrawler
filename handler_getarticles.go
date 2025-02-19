package main

import "net/http"

func (apicfg *apiConfig) handlerGetArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dbArticles, err := apicfg.db.GetRandomFiveArticle(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't get article", err)
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
