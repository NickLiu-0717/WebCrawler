package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/NickLiu-0717/crawler/internal/auth"
	"github.com/NickLiu-0717/crawler/internal/database"
)

func (apicfg *apiConfig) handlerGetArticles(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no authorization header found", err)
		return
	}

	_, err = auth.ValidateJWT(token, apicfg.secretKey)
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

	totolPages, err := apicfg.getTotalPages(limit)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't get total pages", err)
		return
	}
	apicfg.totalPages = totolPages

	w.Header().Set("Content-Type", "application/json")
	dbArticles, err := apicfg.db.GetLatestArticles(r.Context(), database.GetLatestArticlesParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
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

func (apicfg *apiConfig) getTotalPages(limit int) (int, error) {
	totalCount, err := apicfg.db.GetTotalArticleCount(context.Background())
	if err != nil {
		return 0, err
	}
	totalPages := (int(totalCount) + limit - 1) / limit // 計算總頁數
	return totalPages, nil
}
