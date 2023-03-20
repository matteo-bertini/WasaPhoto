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

func (rt *_router) getFollowing(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Non essendo previsto un RequestBody per la richiesta non controllo se è stato specificato o meno.
	// in caso il RequestBody sia stato specificato verrà semplicemente ignorato

	// Estrazione dell'username dalla query nell'URL e controllo dell'esistenza
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
		// Controllo che l'id passato come Authorization corrisponda ad un profilo autorizzato ad ottenere i followers dell'user il cui
		// username è specificato nella query

		// Estrazione dell' Authorization dall'header

		auth_header := r.Header.Get("Authorization")

		// Non è stato specificato il campo Authorization nell'header
		if auth_header == "" {
			w.WriteHeader(http.StatusUnauthorized)
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
					w.WriteHeader(http.StatusUnauthorized)
					return
				} else {
					// Id specificato nel campo Authorization in modo corretto

					// Controllo che l'Id passato come Authorization corrisponda ad un utente esistente ed in tal caso estraggo username
					username1, err := rt.db.UsernameFromId(id1)
					if err != nil {
						if errors.Is(err, utils.ErrUserDoesNotExist) {
							w.WriteHeader(http.StatusUnauthorized)
							ctx.Logger.WithError(err).Error("L'id passato nell' Authorization non corrisponde ad un user registrato.")
							return
						} else {
							w.WriteHeader(http.StatusInternalServerError)
							ctx.Logger.WithError(err).Error("Si è verificato un errore nell'estrazione dell'username dell'Id passato nell'Authorization.")
							return
						}

					} else {
						// Controllo che l'id1 corrisponda ad un profilo esistente
						err = rt.db.CheckUserExistence(*username1)
						if err != nil {
							if errors.Is(err, utils.ErrUserDoesNotExist) {
								w.WriteHeader(http.StatusForbidden)
								ctx.Logger.WithError(err).Error("L'id passato nell' Authorization non corrisponde ad un user con profilo esistente.")
								return

							} else {
								w.WriteHeader(http.StatusInternalServerError)
								ctx.Logger.WithError(err).Error("Si è verificato un errore nell controllare l'esistenza dell'username dell'id passato nell'Authorization.")
								return
							}

						} else {
							// id1 corrisponde ad un user esistente
							// Controllo che id1 sia autorizzato a vedere le informazioni di urlusername
							urlid, err := rt.db.IdFromUsername(urlusername)
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								ctx.Logger.WithError(err).Error("Non è stato possibile estrarre l'Id dell'username nella query.")
								return

							} else {
								err = rt.db.IsAllowed(id1, *urlid)
								if err != nil {
									if errors.Is(err, utils.ErrUserNotAllowed) {
										w.WriteHeader(http.StatusForbidden)
										ctx.Logger.WithError(err).Error("L'user il cui id è passato nell' Authorization non è autorizzato a vedere le informazioni dell'user specificato nella query.")
										return

									} else {
										w.WriteHeader(http.StatusInternalServerError)
										ctx.Logger.WithError(err).Error("Non è stato possibile controllare se l'user il cui id è passato nell' Authorization è autorizzato a vedere le informazioni.")
										return
									}
								} else {
									following, err := rt.db.GetFollowing(*urlid)
									if err != nil {
										w.WriteHeader(http.StatusInternalServerError)
										ctx.Logger.WithError(err).Error("Si è verificato un errore nelle operazioni di database.")
										return

									} else {
										getFollowersResponseBody := getFollowingResponseBody{Following: *following}
										w.Header().Set("Content-Type", "application/json")
										w.WriteHeader(http.StatusOK)
										err = json.NewEncoder(w).Encode(getFollowersResponseBody)
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
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
	}

}
