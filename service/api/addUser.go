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

func (rt *_router) addUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Controllo che la richiesta abbia specificato il RequestBody ed in tal caso lo estraggo
	var addUserRequestBody addUserRequestBody
	err := json.NewDecoder(r.Body).Decode(&addUserRequestBody)
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
		if len(addUserRequestBody.Username) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Error("Il json nel RequestBody presenta degli errori.")
			return

		} else { // Il RequestBody passato non presenta errori

			// Controllo che l'Username passato nel RequestBody sia una stringa conforme alle specifiche
			if utils.CheckUsername(addUserRequestBody.Username) == false {
				w.WriteHeader(http.StatusBadRequest)
				ctx.Logger.Error("L'Username passato nel RequestBody non è conforme alle specifiche.")
				return
			}
			err = rt.db.CheckAuthorization(r, addUserRequestBody.Username)
			if err != nil {
				// L'id non è stato specificato correttamente nell'authorization
				if errors.Is(err, utils.ErrAuthorizationNotSpecified) || errors.Is(err, utils.ErrBearerTokenNotSpecifiedWell) {
					ctx.Logger.WithError(err).Error("Il campo Authorization nell'header presenta degli errori.")
					w.WriteHeader(http.StatusUnauthorized)
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
				id := strings.Split(r.Header.Get("Authorization"), " ")[1]
				err = rt.db.AddUser(addUserRequestBody.Username, id)
				if err != nil {
					if errors.Is(err, utils.ErrUserAlreadyExists) {
						ctx.Logger.WithError(err).Error("L'username passato nel RequestBody corrisponde ad un profilo già esistente.")
						w.WriteHeader(http.StatusForbidden)
						return
					} else {
						w.WriteHeader(http.StatusInternalServerError)
						ctx.Logger.WithError(err).Error("Si è verificato un errore nelle operazioni sul database.")
						return
					}
				} else {
					// L'utente è stato creato,costruzione della risposta
					var addUserResponseBody addUserResponseBody
					addUserResponseBody.Username = addUserRequestBody.Username
					addUserResponseBody.Followers = 0
					addUserResponseBody.Following = 0
					addUserResponseBody.NumberOfPhotos = 0
					addUserResponseBody.UploadedPhotos = []Photo{}
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)
					err = json.NewEncoder(w).Encode(addUserResponseBody)
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
