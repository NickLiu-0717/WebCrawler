package handler

import (
	"net/http"

	"github.com/NickLiu-0717/crawler/internal/auth"
)

type AccessToken struct {
	Token string `json:"token"`
}

func (apicfg *Handler) HandlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no refresh token found", err)
		return
	}

	userID, err := apicfg.Config.Db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token", err)
		return
	}

	token, err := auth.MakeJWT(userID, apicfg.Config.SecretKey)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't create JWT", err)
		return
	}
	accessToken := AccessToken{Token: token}
	respondWithJSON(w, http.StatusOK, accessToken)
}

func (apicfg *Handler) HandlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no refresh token found", err)
	}

	if err = apicfg.Config.Db.UpdateRefreshToken(r.Context(), refreshToken); err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token", err)
	}

	w.WriteHeader(http.StatusNoContent)
}
