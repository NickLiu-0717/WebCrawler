package handler

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/NickLiu-0717/crawler/internal/auth"
	"github.com/NickLiu-0717/crawler/internal/database"
	"github.com/NickLiu-0717/crawler/internal/models"
)

func (apicfg *Handler) HandlerGetArticles(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no authorization header found", err)
		return
	}

	_, err = auth.ValidateJWT(token, apicfg.Config.SecretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired access token", err)
		return
	}

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

	totolPages, err := apicfg.GetTotalPages(limit)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't get total pages", err)
		return
	}
	apicfg.Config.TotalPages = totolPages

	w.Header().Set("Content-Type", "application/json")
	dbArticles, err := apicfg.Config.Db.GetLatestArticles(r.Context(), database.GetLatestArticlesParams{
		Limit:  lim,
		Offset: off,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't get article", err)
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

func (apicfg *Handler) GetTotalPages(limit int) (int, error) {
	totalCount, err := apicfg.Config.Db.GetTotalArticleCount(context.Background())
	if err != nil {
		return 0, err
	}
	totalPages := (int(totalCount) + limit - 1) / limit
	return totalPages, nil
}

func safeIntTo32(i int) (int32, error) {
	if i > math.MaxInt32 || i < math.MinInt32 {
		return 0, fmt.Errorf("value %d is out of int32 range", i)
	}
	return int32(i), nil
}
