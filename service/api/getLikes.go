package api

import (
	"WasaPhoto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoid := strings.Split(r.URL.Path, "/")[4]
	likes, err := rt.db.GetLikes(photoid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Si è verificato un errore nelle operazioni di database.")
		return

	} else {
		getLikesResponseBody := getLikesResponseBody{Likes: *likes}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(getLikesResponseBody)
		// Si è verificato un errore nell'encoding della risposta
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Si è verificato un errore nell'encoding della risposta.")
			return
		} else {
			// Non si sono verificati errori,ritorno
			return
		}

	}

}
