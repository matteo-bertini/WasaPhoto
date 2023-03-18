package api

import (
	"WasaPhoto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoid := strings.Split(r.URL.Path, "/")[4]
	comments, err := rt.db.GetComments(photoid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("Si è verificato un errore nelle operazioni di database.")
		return

	} else {
		getCommentsResponseBody := getCommentsResponseBody{Comments: *comments}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(getCommentsResponseBody)
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
