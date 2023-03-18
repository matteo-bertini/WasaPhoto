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
	"github.com/segmentio/ksuid"
)

func (rt *_router) commentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	} else {
		// L'user nell'URL esiste
		urlid, err := rt.db.IdFromUsername(urlusername)
		if err != nil {
			// Si è verificato un problema nell'estrarre l'Id
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Si è verificato un problema nell'estrarre l'Id.")
			return
		} else {
			photoid := strings.Split(r.URL.Path, "/")[4]
			err = rt.db.CheckPhotoExistence(*urlid, photoid)
			if err != nil {
				if errors.Is(err, utils.ErrPhotoDoesNotExist) {
					w.WriteHeader(http.StatusNotFound)
					ctx.Logger.WithError(err).Error("La foto nell'URL non esiste.")
					return

				} else {
					w.WriteHeader(http.StatusInternalServerError)
					ctx.Logger.WithError(err).Error("Si è verificato un errore nel controllare l'esistenza della foto nel database.")
					return

				}
			} else {
				// La foto nell'URL esiste,controllo che sia stato specificato il RequestBody

				var commentPhotoRequestBody commentPhotoRequestBody
				err = json.NewDecoder(r.Body).Decode(&commentPhotoRequestBody)
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
					if len(commentPhotoRequestBody.CommentAuthor) == 0 || len(commentPhotoRequestBody.CommentText) == 0 {
						w.WriteHeader(http.StatusBadRequest)
						ctx.Logger.Error("Il json nel RequestBody presenta degli errori.")
						return

					} else { // Il RequestBody passato non presenta errori
						// Controllo che l'autore del commento esista e che sia autorizzato
						err := rt.db.CheckUserExistence(commentPhotoRequestBody.CommentAuthor)
						if err != nil {
							// L'user non esiste
							if errors.Is(err, utils.ErrUserDoesNotExist) {
								w.WriteHeader(http.StatusNotFound)
								ctx.Logger.WithError(err).Error("L'autore del commento non è un user esistente.")
								return
							} else {
								// Si è verificato un problema nel controllare l'esistenza dell'user
								w.WriteHeader(http.StatusInternalServerError)
								ctx.Logger.WithError(err).Error("Si è verificato un errore nel controllare l'esistenza dell'autore del commento nel database.")
								return
							}

						} else {
							err = rt.db.CheckAuthorization(r, commentPhotoRequestBody.CommentAuthor)
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
								// Controllo che l'user sia autorizzato a commentare
								err = rt.db.IsAllowed(strings.Split(r.Header.Get("Authorization"), " ")[1], *urlid)
								if err != nil {
									if errors.Is(err, utils.ErrUserNotAllowed) {
										w.WriteHeader(http.StatusUnauthorized)
										ctx.Logger.WithError(err).Error("L'user non è autorizzato a commentare.")
										return
									} else {
										w.WriteHeader(http.StatusInternalServerError)
										ctx.Logger.WithError(err).Error("Si è verificato un errore nel controllare i ban.")
										return

									}
								} else {
									// Generazione del CommentId
									commentPhotoRequestBody.CommentId = ksuid.New().String()
									err = rt.db.CommentPhoto(*urlid, photoid, commentPhotoRequestBody.CommentId, strings.Split(r.Header.Get("Authorization"), " ")[1], commentPhotoRequestBody.CommentText)
									if err != nil {
										w.WriteHeader(http.StatusInternalServerError)
										ctx.Logger.WithError(err).Error("Si è verificato un errore nelle operazioni di database.")
										return
									} else {
										w.Header().Set("Content-Type", "application/json")
										w.WriteHeader(http.StatusCreated)
										err = json.NewEncoder(w).Encode(commentPhotoRequestBody)
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

	}
}
