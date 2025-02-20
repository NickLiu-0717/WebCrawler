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
	}

	hashPassword, err := auth.HashPassword(userInfo.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't hash password", err)
	}
	dbUser, err := apicfg.db.CreateNewUser(r.Context(), database.CreateNewUserParams{
		Email:    userInfo.Email,
		Password: hashPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't create new user", err)
	}
	user := User{
		ID:         dbUser.ID,
		Created_at: dbUser.CreatedAt,
		Updated_at: dbUser.UpdatedAt,
		Email:      dbUser.Email,
		Password:   dbUser.Password,
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (apicfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var userInfo UserInfo
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't decode request body", err)
	}

}
