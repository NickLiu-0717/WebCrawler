package main

import (
	"encoding/json"
	"net/http"

	"github.com/NickLiu-0717/crawler/internal/auth"
	"github.com/NickLiu-0717/crawler/internal/database"
)

func (apicfg *apiConfig) handlerSignup(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var userInfo UserInfo

	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode request body", err)
		return
	}

	hashPassword, err := auth.HashPassword(userInfo.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't hash password", err)
		return
	}
	dbUser, err := apicfg.db.CreateNewUser(r.Context(), database.CreateNewUserParams{
		Email:    userInfo.Email,
		Password: hashPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't create new user", err)
		return
	}
	user := User{
		ID:         dbUser.ID,
		Created_at: dbUser.CreatedAt,
		Updated_at: dbUser.UpdatedAt,
		Email:      dbUser.Email,
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (apicfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var userInfo UserInfo
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode request body", err)
		return
	}

	dbUser, err := apicfg.db.GetUserFromEmail(r.Context(), userInfo.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "user not found", err)
		return
	}

	if err = auth.CheckPasswordHash(userInfo.Password, dbUser.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect password", err)
		return
	}

	token, err := auth.MakeJWT(dbUser.ID, apicfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't make JWT token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't make refresh token", err)
		return
	}

	_, err = apicfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: dbUser.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't create refresh token", err)
		return
	}

	user := User{
		ID:           dbUser.ID,
		Created_at:   dbUser.CreatedAt,
		Updated_at:   dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        token,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, user)
}
