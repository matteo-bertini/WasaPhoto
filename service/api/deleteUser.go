package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/utils"
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// La richiesta non ha body,quindi non controllo se sia stato specificato o meno.
	// In caso sia stato specificato lo ignoro.

	// Estraggo l'username dall'URL e ne controllo l'esistenza.
	username := strings.Split(r.URL.Path, "/")[2]
	err := rt.db.CheckUserExistence(username)
	if err != nil {
		if errors.Is(err, utils.ErrUserDoesNotExist) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("L'username specificato nell'URL non esiste.")
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Errore nella verifica dell'esistenza dell' username.")
			return
		}
	} else {
		// l'username specificato nell'URL corrispondono a user esistenti
		err = rt.db.CheckAuthorization(r, username)
		if err != nil {
			// L'id non è stato specificato correttamente nell'authorization
			if errors.Is(err, utils.ErrorAuthorizationNotSpecified) || errors.Is(err, utils.ErrorBearerTokenNotSpecifiedWell) {
				ctx.Logger.WithError(err).Error("Il campo Authorization nell'header presenta degli errori.")
				w.WriteHeader(http.StatusBadRequest)
				return
				// L'id non è autorizzato ad effettuare l'operazione
			} else if errors.Is(err, utils.ErrorUnauthorized) {
				ctx.Logger.WithError(err).Error("L'id passato non è autorizzato ad effettuare l'operazione.")
				w.WriteHeader(http.StatusUnauthorized)
				return

			} else {
				// Errore nell'esecuzione delle query
				ctx.Logger.WithError(err).Error("Si è verificato un errore nella verifica dell'id all'interno del database.")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			// Autorizzato
			err = rt.db.DeleteUser(strings.Split(r.Header.Get("Authorization"), " ")[1], username)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("Si è verificato un errore nelle operazioni di database.")
			} else {
				w.WriteHeader(http.StatusNoContent)
				return
			}

		}

	}

}
