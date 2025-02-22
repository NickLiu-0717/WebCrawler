package main

import "net/http"

func (apicfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if err := apicfg.db.DeleteArticles(r.Context()); err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't delete articles", err)
		return
	}

	if err := apicfg.db.DeleteUsers(r.Context()); err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't delete users", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset Successfully"))
}
