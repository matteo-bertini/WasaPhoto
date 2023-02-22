package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estraggo il vecchio username dall' URL
	old_username := strings.Split(r.URL.Path, "/")[2]
	err := rt.db.CheckUserExistence(old_username)
	if err != nil {
		if errors.Is(err, utils.ErrUserDoesNotExist) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("L' username specificato nell'URL non corrisponde ad un user esistente.")
			return

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Si è verificato un errore nel database nel controllare l'esistenza dell'user.")

		}
	} else {
		// L'username specificato nell'URL corrisponde da un user esistente

		// Controllo che la richiesta abbia specificato il RequestBody ed in tal caso lo estraggo
		var setMyUsernameRequestBody setMyUsernameRequestBody
		err := json.NewDecoder(r.Body).Decode(&setMyUsernameRequestBody)
		if err != nil {
			// Non è stato specificato il RequestBody per la richiesta
			if errors.Is(err, io.EOF) {
				ctx.Logger.WithError(err).Error("Non è stato specificato il RequestBody per la richiesta.")
				w.WriteHeader(http.StatusBadRequest)
				return

			} else {
				// Il RequestBody è stato passato,ma non è stato possibile decodificarlo
				ctx.Logger.WithError(err).Error("Non è stato possibile decodificare il RequestBody.")
				w.WriteHeader(http.StatusBadRequest)
				return

			}
		} else { // Il RequestBody è stato decodficato

			// C'è un errore nel RequestBody passato (nomi dei campi errati,campi non specificati,ecc)
			if len(setMyUsernameRequestBody.Username) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				ctx.Logger.Error("Il json nel RequestBody presenta degli errori.")
				return

			} else { // Il RequestBody passato non presenta errori

				// Controllo che l'Username passato nel RequestBody sia una stringa conforme alle specifiche
				if utils.CheckUsername(setMyUsernameRequestBody.Username) == false {
					w.WriteHeader(http.StatusBadRequest)
					ctx.Logger.Error("L'Username passato nel RequestBody non è conforme alle specifiche.")
					return
				} else {
					// L'username passato nel RequestBody è conforme
					err = rt.db.CheckAuthorization(r, old_username)
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
						err = rt.db.SetMyUsername(old_username, setMyUsernameRequestBody.Username)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							ctx.Logger.WithError(err).Error("Si è verificato un errore nelle operazioni di database.")
							return
						} else {
							w.Header().Set("Content-Type", "application/json")
							w.WriteHeader(http.StatusOK)
							err = json.NewEncoder(w).Encode(setMyUsernameRequestBody)
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

		}
	}
}
