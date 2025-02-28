package handler

import "net/http"

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
	w.Write([]byte("Reset Successfully"))
}
