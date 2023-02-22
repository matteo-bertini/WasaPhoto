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

func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo l'esistenza dell'user specificato nell'URL
	urlusername := strings.Split(r.URL.Path, "/")[2]
	err := rt.db.CheckUserExistence(urlusername)
	if err != nil {
		if errors.Is(err, utils.ErrUserDoesNotExist) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("L'user specificato nell'URL non esiste.")
			return

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Si è verificato un problema nel controllare l'esistenza dell'user nell'URL.")
			return
		}
	} else {
		// L'user nell'URL esiste
		// Controllo che sia stato specificato Authorization
		auth_header := r.Header.Get("Authorization")

		// Non è stato specificato il campo Authorization nell'header
		if auth_header == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			// Il campo Authorization è stato specificato nell'header (Authorization : Bearer abcdef)
			splitted_authorization := strings.Split(auth_header, " ")

			// Controllo che l'Id sia stato specificato correttamente nel campo Authorization
			if len(splitted_authorization) == 2 {
				authorization_type := splitted_authorization[0]
				id1 := splitted_authorization[1]
				// Id non specificato in conformità con le specifiche
				if authorization_type != utils.Bearer_Authorization || strings.TrimSpace(id1) == "" {
					w.WriteHeader(http.StatusBadRequest)
					return
				} else {
					// Id specificato nel campo Authorization in modo corretto
					id1 := strings.Split(r.Header.Get("Authorization"), " ")[1]
					username1, err := rt.db.UsernameFromId(id1)
					if err != nil {
						if errors.Is(err, utils.ErrUserDoesNotExist) {
							w.WriteHeader(http.StatusNotFound)
							ctx.Logger.WithError(err).Error("L'user il cui id è passato nell'authorization non è registrato.")
						} else {
							w.WriteHeader(http.StatusInternalServerError)
							ctx.Logger.WithError(err).Error("Si è verificato un errore nell'estrarre l'username dell'id nell'Authorization.")
							return
						}
					} else {
						err = rt.db.CheckUserExistence(*username1)
						if err != nil {
							if errors.Is(err, utils.ErrUserDoesNotExist) {
								w.WriteHeader(http.StatusNotFound)
								ctx.Logger.WithError(err).Error("L'user specificato nell'URL non esiste.")
								return

							} else {
								w.WriteHeader(http.StatusInternalServerError)
								ctx.Logger.WithError(err).Error("Si è verificato un problema nel controllare l'esistenza dell'user nell'URL.")
								return
							}
						} else {
							// entrambi gli user esistono
							id2, err := rt.db.IdFromUsername(urlusername)
							if err != nil {
								// Controllo solo un caso,l'user esiste già
								w.WriteHeader(http.StatusInternalServerError)
								ctx.Logger.WithError(err).Error("Si è verificato un problema nell'estrazione dell'Id dell'user nell'URL.")
								return
							} else {
								err = rt.db.IsAllowed(id1, *id2)
								if err != nil {
									if errors.Is(err, utils.ErrUserNotAllowed) {
										w.WriteHeader(http.StatusUnauthorized)
										ctx.Logger.WithError(err).Error("L'user non è autorizzato a mettere like.")
										return
									} else {
										w.WriteHeader(http.StatusInternalServerError)
										ctx.Logger.WithError(err).Error("Si è verficato un errore nel controlalre l'autorizzazione.")
										return

									}
								} else {
									err = rt.db.CheckPhotoExistence(*id2, strings.Split(r.URL.Path, "/")[4])
									if err != nil {
										if errors.Is(err, utils.ErrPhotoDoesNotExist) {
											w.WriteHeader(http.StatusNotFound)
											ctx.Logger.WithError(err).Error("La foto cercata non esiste.")
											return

										} else {
											w.WriteHeader(http.StatusInternalServerError)
											ctx.Logger.WithError(err).Error("Si è verficato un errore nel controlalre l'esistenza della foto.")
											return

										}
									}
									err = rt.db.LikePhoto(*id2, strings.Split(r.URL.Path, "/")[4], id1)
									if err != nil {
										w.WriteHeader(http.StatusInternalServerError)
										ctx.Logger.WithError(err).Error("Si è verificato un errore nelle operazioni sul database.")
										return

									} else {
										var likePhotoResponseBody likePhotoResponseBody
										likePhotoResponseBody.LikeId = strings.Split(r.URL.Path, "/")[4]
										w.Header().Set("Content-Type", "application/json")
										w.WriteHeader(http.StatusCreated)
										err = json.NewEncoder(w).Encode(likePhotoResponseBody)
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
			} else {
				w.WriteHeader(http.StatusBadRequest)
				ctx.Logger.Error("La richiesta non ha specificato correttamente l'Authorization.")
				return
			}
		}

	}
}
