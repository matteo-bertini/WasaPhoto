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

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estrazione dell'username dall' URL e controllo dell'esistenza
	urlusername := strings.Split(r.URL.Path, "/")[2]
	err := rt.db.CheckUserExistence(urlusername)
	if err != nil {
		// L'user non esiste
		if errors.Is(err, utils.ErrUserDoesNotExist) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("L'user specificato nell' URL non esiste.")
			return
		} else {
			// Si è verificato un problema nel controllare l'esistenza dell'user
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Si è verificato un errore nel controllare l'esistenza dell'user nel database.")
			return
		}
	}
	// L'user specificato nell'URL esiste
	// Controllo che la richiesta abbia specificato il RequestBody ed in tal caso lo estraggo
	var followUserRequestBody followUserRequestBody
	err = json.NewDecoder(r.Body).Decode(&followUserRequestBody)
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
		if len(followUserRequestBody.FollowerId) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Error("Il json nel RequestBody presenta degli errori.")
			return

		} else { // Il RequestBody passato non presenta errori

			// Controllo che l'Username passato nel RequestBody sia una stringa conforme alle specifiche
			if utils.CheckUsername(followUserRequestBody.FollowerId) == false {
				w.WriteHeader(http.StatusBadRequest)
				ctx.Logger.Error("L'Username passato nel RequestBody non è conforme alle specifiche.")
				return
			}
			err = rt.db.CheckAuthorization(r, followUserRequestBody.FollowerId)
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

				// Controllo che l'username specificato nel RequestBody corrisponda ad un user esistente
				err := rt.db.CheckUserExistence(followUserRequestBody.FollowerId)
				if err != nil {
					// L'user non esiste
					if errors.Is(err, utils.ErrUserDoesNotExist) {
						w.WriteHeader(http.StatusForbidden)
						ctx.Logger.WithError(err).Error("L'operazione non può essere completata,l'user specificato nel RequestBody non esiste.")
						return
					} else {
						// Si è verificato un problema nel controllare l'esistenza dell'user
						w.WriteHeader(http.StatusInternalServerError)
						ctx.Logger.WithError(err).Error("Si è verificato un errore nel controllare l'esistenza dell'user nel database.")
						return
					}

				} else {
					// L'username specificato nel RequestBody corrisponde ad un user esistente
					url_id, err := rt.db.IdFromUsername(urlusername)
					if err != nil {
						if errors.Is(err, utils.ErrUserDoesNotExist) {
							w.WriteHeader(http.StatusNotFound)
							ctx.Logger.WithError(err).Error("L'user specificato nell'URL non esiste.")
							return

						} else {
							w.WriteHeader(http.StatusInternalServerError)
							ctx.Logger.WithError(err).Error("Si è verificato un errore nell'estrarre l'Id dell'username nell' URL.")
							return
						}
					} else {
						if urlusername == followUserRequestBody.FollowerId {
							w.WriteHeader(http.StatusForbidden)
							ctx.Logger.Error("Un user non può seguire se stesso.")

						} else {
							err = rt.db.FollowUser(followUserRequestBody.FollowerId, strings.Split(r.Header.Get("Authorization"), " ")[1], urlusername, *url_id)
							if err != nil {
								if errors.Is(err, utils.ErrFollowerAlreadyAdded) {
									w.WriteHeader(http.StatusForbidden)
									ctx.Logger.WithError(err).Error("L'user è già presente nella lista dei followers dell'altro.")
									return
								}
								w.WriteHeader(http.StatusInternalServerError)
								ctx.Logger.WithError(err).Error("Si è verificato un errore nelle operazioni sul database.")
								return
							} else {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusCreated)
								err = json.NewEncoder(w).Encode(followUserRequestBody)
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

}
