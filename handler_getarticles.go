package main

import "net/http"

func (apicfg *apiConfig) handlerGetArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	dbArticle, err := apicfg.db.GetOneArticle(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't get article", err)
		return
	}

	article := Article{
		ID:           dbArticle.ID,
		URL:          dbArticle.Url,
		Title:        dbArticle.Title,
		Content:      dbArticle.Content,
		Catagory:     dbArticle.Catagory,
		Created_at:   dbArticle.CreatedAt,
		Published_at: dbArticle.PublishedAt,
	}
	respondWithJSON(w, http.StatusOK, article)
}
