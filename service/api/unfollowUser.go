package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/utils"
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// La richiesta non ha body,quindi non controllo se sia stato specificato o meno.
	// In caso sia stato specificato lo ignoro.

	// Estraggo i due username dall'URL e controllo l'esistenza di entrambi.
	username := strings.Split(r.URL.Path, "/")[2]
	to_del := strings.Split(r.URL.Path, "/")[4]
	err := rt.db.CheckUserExistence(username)
	if err != nil {
		if errors.Is(err, utils.ErrUserDoesNotExist) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("Il primo username specificato nell'URL non esiste.")
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Errore nella verifica dell'esistenza del primo username.")
			return
		}
	} else {
		err = rt.db.CheckUserExistence(to_del)
		if err != nil {
			if errors.Is(err, utils.ErrUserDoesNotExist) {
				w.WriteHeader(http.StatusNotFound)
				ctx.Logger.WithError(err).Error("Il secondo username specificato nell'URL non esiste.")
				return
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				ctx.Logger.WithError(err).Error("Errore nella verifica dell'esistenza del secondo username.")
				return
			}

		} else {
			// Entrambi gli username specificati nell'URL corrispondono a user esistenti

			// Controllo che i due username non coincidano,un user non può unfolloware se stesso
			if username == to_del {
				w.WriteHeader(http.StatusNoContent)
				ctx.Logger.Error("Un user non può togliere il follow a se stesso.")
				return
			}
			err = rt.db.CheckAuthorization(r, to_del)
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
			} else {
				// Autorizzato
				id1, err := rt.db.IdFromUsername(username)
				if err != nil {
					if errors.Is(err, utils.ErrUserDoesNotExist) {
						w.WriteHeader(http.StatusNotFound)
						ctx.Logger.WithError(err).Error("Il primo username specificato nell'URL non esiste.")
						return
					} else {
						w.WriteHeader(http.StatusInternalServerError)
						ctx.Logger.WithError(err).Error("Errore nell'estrazione dell' id del primo unsername.")
						return
					}
				} else {
					err = rt.db.UnfollowUser(username, *id1, to_del, strings.Split(r.Header.Get("Authorization"), " ")[1])
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						ctx.Logger.WithError(err).Error("Si è verificato un errore nelle operazioni sul database.")
						return
					} else {
						w.WriteHeader(http.StatusNoContent)
						return
					}

				}

			}

		}
	}
}
