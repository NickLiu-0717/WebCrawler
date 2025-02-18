package main

import "net/http"

func (apicfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	err := apicfg.db.DeleteArticles(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't delete articles", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Reset Successfully"))
}
