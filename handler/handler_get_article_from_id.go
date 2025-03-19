package handler

import (
	"net/http"

	"github.com/NickLiu-0717/crawler/internal/models"
	"github.com/google/uuid"
)

func (apicfg *Handler) HandlerGetArticleFromID(w http.ResponseWriter, r *http.Request) {
	aID := r.PathValue("articleId")
	articleID, err := uuid.Parse(aID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't parse uuid", err)
		return
	}

	dbArticle, err := apicfg.Config.Db.GetArticleByID(r.Context(), articleID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "article not found", err)
		return
	}

	article := models.Article{
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
