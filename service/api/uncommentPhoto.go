package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/utils"
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) uncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// La richiesta non ha body,quindi non controllo se sia stato specificato o meno.
	// In caso sia stato specificato lo ignoro.

	// Controllo che l'username nell'URL corrisponda ad un user esistente
	urlusername := strings.Split(r.URL.Path, "/")[2]
	err := rt.db.CheckUserExistence(urlusername)
	if err != nil {
		if errors.Is(err, utils.ErrUserDoesNotExist) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("L'user nell'URL non esiste.")
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Si è verificato un errore nel verificare l'esistenza dell'user nell'URL.")
			return

		}
	} else {
		// L'user nell'URL esiste,controllo che esista la foto specificata
		photoid := strings.Split(r.URL.Path, "/")[4]

		// estraggo l'id dell'user nell'URL per controllare l'esistenza della foto
		urlid, err := rt.db.IdFromUsername(urlusername)
		if err != nil {
			// non controllo il caso in cui l'user non esiste
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Si è verificato un errore nell'estrazione dell'Id dell'user nell'URL.")
			return

		} else {
			// Check esistenza della foto
			err = rt.db.CheckPhotoExistence(*urlid, photoid)
			if err != nil {
				if errors.Is(err, utils.ErrPhotoDoesNotExist) {
					w.WriteHeader(http.StatusNotFound)
					ctx.Logger.WithError(err).Error("La foto non esiste.")
					return
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					ctx.Logger.WithError(err).Error("Si è verificato un errore nel controllare l'esistenza della foto.")
					return

				}
			} else {
				// La foto esiste
				// Devo controllare  l'Id passato nell'Authorization
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
							// Controllo che l'Id passato nell'Authorization sia corrispondente ad un utente esistente ed autorizzato a vedere la foto
							username, err := rt.db.UsernameFromId(strings.Split(r.Header.Get("Authorization"), " ")[1])
							if err != nil {
								if errors.Is(err, utils.ErrUserDoesNotExist) {
									w.WriteHeader(http.StatusUnauthorized)
									ctx.Logger.WithError(err).Error("L'authorization non corrisponde ad un user registrato.")
									return
								} else {
									w.WriteHeader(http.StatusInternalServerError)
									ctx.Logger.WithError(err).Error("Errore nell'estrarre username dall'Id")
									return

								}
							} else {
								err = rt.db.CheckUserExistence(*username)
								if err != nil {
									if errors.Is(err, utils.ErrUserDoesNotExist) {
										w.WriteHeader(http.StatusUnauthorized)
										ctx.Logger.WithError(err).Error("L'authorization non corrisponde ad un user con profilo esistente.")
										return
									} else {
										w.WriteHeader(http.StatusInternalServerError)
										ctx.Logger.WithError(err).Error("Errore nel controllare l'esistenza dell'user.")
										return

									}

								} else {
									err = rt.db.IsAllowed(strings.Split(r.Header.Get("Authorization"), " ")[1], *urlid)
									if err != nil {
										if errors.Is(err, utils.ErrUserNotAllowed) {
											w.WriteHeader(http.StatusUnauthorized)
											ctx.Logger.WithError(err).Error("L'user non è autorizzato a vedere le informazioni .")
											return
										} else {
											w.WriteHeader(http.StatusInternalServerError)
											ctx.Logger.WithError(err).Error("Errore nel controllare l'autorizzazione.")
											return

										}
									} else {
										err = rt.db.UncommentPhoto(*urlid, photoid, strings.Split(r.URL.Path, "/")[6], strings.Split(r.Header.Get("Authorization"), " ")[1])
										if err != nil {
											w.WriteHeader(http.StatusInternalServerError)
											ctx.Logger.WithError(err).Error("Si è verificato un errore nelle operazioni di database.")
											return
										} else {
											w.WriteHeader(http.StatusNoContent)
											return
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

	}

}
