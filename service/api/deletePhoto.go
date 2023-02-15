package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/utils"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo che l'username nell'URL corrisponde ad un profilo esistente
	urlusername := strings.Split(r.URL.Path, "/")[2]
	err := rt.db.CheckUserExistence(urlusername)
	if err != nil {
		if errors.Is(err, utils.ErrUserDoesNotExist) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("L'username nell'url non corrisponde ad un user esistente.")
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Si è verificato un errore nel controllare l'esistenza dell'username nell'URL.")
			return

		}
	} else {
		// L'user esiste
		err = rt.db.CheckAuthorization(r, urlusername)
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

			// Prima elimino la foto dal server poi dal database
			photoid := strings.Split(r.URL.Path, "/")[4]
			err = os.Remove("/tmp/WasaPhoto/" + strings.Split(r.Header.Get("Authorization"), " ")[1] + "/" + photoid + ".png")
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				ctx.Logger.WithError(err).Error("La foto cercata non esiste.")
				return
			} else {
				// Rimuovo la foto dal database
				err = rt.db.DeletePhoto(strings.Split(r.Header.Get("Authorization"), " ")[1], photoid)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					ctx.Logger.WithError(err).Error("Si è verificato un errore nelle operazioni di database")
					return
				} else {
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}
		}

	}

}
