package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// La richiesta non ha RequestBody,quindi ne ignoro il controllo

	// Controllo l'esistenza dell'username passato nell' URL
	urlusername := strings.Split(r.URL.Path, "/")[2]
	err := rt.db.CheckUserExistence(urlusername)
	if err != nil {
		if errors.Is(err, utils.ErrUserDoesNotExist) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("L'username nell'URL non corrisponde ad un user esistente.")
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Si è verificato un errore nel controllare l'esistenza dell'username nell'URL.")
			return
		}
	} else {
		// L'username nell'URL corrisponde ad un user esistente
		err = rt.db.CheckAuthorization(r, urlusername)
		if err != nil {
			// L'id non è stato specificato correttamente nell'authorization
			if errors.Is(err, utils.ErrAuthorizationNotSpecified) || errors.Is(err, utils.ErrBearerTokenNotSpecifiedWell) {
				ctx.Logger.WithError(err).Error("Il campo Authorization nell'header presenta degli errori.")
				w.WriteHeader(http.StatusBadRequest)
				return
				// L'id non è autorizzato ad effettuare l'operazione
			} else if errors.Is(err, utils.ErrUnauthorized) {
				ctx.Logger.WithError(err).Error("L'id passato non è autorizzato ad effettuare l'operazione.")
				w.WriteHeader(http.StatusUnauthorized)
				return

			} else {
				// Errore nell'esecuzione delle query
				ctx.Logger.WithError(err).Error("Si è verificato un errore nella verifica dell'id all'interno del database.")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		urlid, err := rt.db.IdFromUsername(urlusername)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Si è verificato un errore nell'estrazione dell'Id dell'username nell'URL.")
			return

		} else {
			photostream, err := rt.db.GetMyStream(*urlid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("Si è verificato un errore nelle operazioni di database.")
				return
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				err = json.NewEncoder(w).Encode(photostream)
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

	}

}
