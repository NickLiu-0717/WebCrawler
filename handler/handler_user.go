package handler

import (
	"encoding/json"
	"net/http"

	"github.com/NickLiu-0717/crawler/internal/auth"
	"github.com/NickLiu-0717/crawler/internal/database"
	"github.com/NickLiu-0717/crawler/internal/models"
)

func (apicfg *Handler) HandlerSignup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var userInfo models.UserInfo

	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode request body", err)
		return
	}

	hashPassword, err := auth.HashPassword(userInfo.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't hash password", err)
		return
	}
	dbUser, err := apicfg.Config.Db.CreateNewUser(r.Context(), database.CreateNewUserParams{
		Email:    userInfo.Email,
		Password: hashPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't create new user", err)
		return
	}
	user := models.User{
		ID:         dbUser.ID,
		Created_at: dbUser.CreatedAt,
		Updated_at: dbUser.UpdatedAt,
		Email:      dbUser.Email,
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (apicfg *Handler) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var userInfo models.UserInfo
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode request body", err)
		return
	}

	dbUser, err := apicfg.Config.Db.GetUserFromEmail(r.Context(), userInfo.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user not found", err)
		return
	}

	if err = auth.CheckPasswordHash(userInfo.Password, dbUser.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect password", err)
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, apicfg.Config.SecretKey)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't make JWT token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't make refresh token", err)
		return
	}

	_, err = apicfg.Config.Db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: dbUser.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't create refresh token", err)
		return
	}

	user := models.User{
		ID:           dbUser.ID,
		Created_at:   dbUser.CreatedAt,
		Updated_at:   dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, user)
}
