package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (apicfg *apiConfig) handlerGetArticleFromID(w http.ResponseWriter, r *http.Request) {
	aID := r.PathValue("articleId")
	articleID, err := uuid.Parse(aID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't parse uuid", err)
		return
	}

	dbArticle, err := apicfg.db.GetArticleByID(r.Context(), articleID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "article not found", err)
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
