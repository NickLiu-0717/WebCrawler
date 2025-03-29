package handler

import (
	"net/http"
	"strconv"

	"github.com/NickLiu-0717/crawler/internal/database"
	"github.com/NickLiu-0717/crawler/internal/models"
)

func (apicfg *Handler) HandlerGetCategoryArticles(w http.ResponseWriter, r *http.Request) {
	cate := r.PathValue("category")

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 5
	}

	offset := (page - 1) * limit

	lim, err := safeIntTo32(limit)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "int overflow", err)
		return
	}
	off, err := safeIntTo32(offset)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "int overflow", err)
		return
	}

	dbArticles, err := apicfg.Config.Db.GetLatestCategoryArticles(r.Context(), database.GetLatestCategoryArticlesParams{
		Catagory: cate,
		Limit:    lim,
		Offset:   off,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't get articles by category from database", err)
		return
	}

	var articles []models.Article
	for _, dbArticle := range dbArticles {
		article := models.Article{
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
