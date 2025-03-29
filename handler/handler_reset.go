package handler

import (
	"log"
	"net/http"
)

func (apicfg *Handler) HandlerReset(w http.ResponseWriter, r *http.Request) {
	if err := apicfg.Config.Db.DeleteArticles(r.Context()); err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't delete articles", err)
		return
	}

	if err := apicfg.Config.Db.DeleteUsers(r.Context()); err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't delete users", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("Reset Successfully"))
	if err != nil {
		log.Printf("Error writing JSON: %s", err)
		w.WriteHeader(500)
		return
	}
}
