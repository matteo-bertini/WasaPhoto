package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/utils"
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Estraggo l'username dall' URL e ne controllo l'esistenza
	urlusername := strings.Split(r.URL.Path, "/")[2]
	err := rt.db.CheckUserExistence(urlusername)
	if err != nil {
		if errors.Is(err, utils.ErrUserDoesNotExist) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("L'username nell'URL non corrisponde da un user con un profilo esistente.")
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Si è verificato un errore nel controllare l'esistenza dell'user nell'URL.")
			return

		}
	} else {
		// L'username nell'URL corrisponde ad un user esistente
		// Controllo che l'id passato come Authorization corrisponda ad un profilo autorizzato a vedere il profilo dell'user il cui
		// username è specificato nella query
		// Estrazione dell' Authorization dall'header

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
					username, err := rt.db.UsernameFromId(id1)
					if err != nil {
						if errors.Is(err, utils.ErrUserDoesNotExist) {
							w.WriteHeader(http.StatusNotFound)
							ctx.Logger.WithError(err).Error("L'Id nell'authorization non corrisponde ad un user registrato.")
							return

						} else {
							w.WriteHeader(http.StatusInternalServerError)
							ctx.Logger.WithError(err).Error("Si è verificato un errore nell'estrazione dell'username dall'Id nel Authorization.")
							return

						}
					} else {
						err = rt.db.CheckUserExistence(*username)
						if err != nil {
							if errors.Is(err, utils.ErrUserDoesNotExist) {
								w.WriteHeader(http.StatusNotFound)
								ctx.Logger.WithError(err).Error("L'Id nell'authorization non corrisponde ad un user con profilo esistente.")
								return

							} else {
								w.WriteHeader(http.StatusInternalServerError)
								ctx.Logger.WithError(err).Error("Si è verificato un errore nel controllare l'esistenza del secondo user.")
								return

							}
						} else {
							// Entrambi gli user esistono
							id2, err := rt.db.IdFromUsername(urlusername)
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								ctx.Logger.WithError(err).Error("Si è verificato un errore nell'estrazione dell'id dell'username nell'URL.")
								return
							}
							err = rt.db.IsAllowed(id1, *id2)
							if err != nil {
								if errors.Is(err, utils.ErrUserNotAllowed) {
									w.WriteHeader(http.StatusUnauthorized)
									ctx.Logger.WithError(err).Error("L'user non è autorizzato a vedere le informazioni.")
									return

								}
								w.WriteHeader(http.StatusInternalServerError)
								ctx.Logger.WithError(err).Error("Si è verificato un errore nel controllare l'autorizzazione per le informazioni.")
								return

							} else {
								// L'user è autorizzato a vedere la foto,costruisco la risposta
								photoid := strings.Split(r.URL.Path, "/")[4]
								path := "/tmp/WasaPhoto/" + *id2 + "/" + photoid + ".png"
								photofile, err := os.Open(path)
								if err != nil {
									ctx.Logger.WithError(err).Error("La foto cercata non esiste.")
									w.WriteHeader(http.StatusNotFound)
									return
								} else {

									// Setting dell'header della risposta
									w.Header().Set("Content-Type", "image/png")

									buf := bytes.NewBuffer(nil)
									_, err := io.Copy(buf, photofile)
									if err != nil {
										w.WriteHeader(http.StatusInternalServerError)
										ctx.Logger.WithError(err).Error("Si è verificato un errore nella copia sul buffer.")
										return
									} else {
										w.WriteHeader(http.StatusOK)
										_, err = w.Write(buf.Bytes())
										if err != nil {
											ctx.Logger.WithError(err).Error("Si è verificato un errore nella scrittura della riposta.")
											w.WriteHeader(http.StatusInternalServerError)
											return
										}
										return
									}
								}
							}
						}

					}

				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				ctx.Logger.Error("La roichiesta non ha specificato correttamente l'Authorization.")
				return
			}
		}
	}
}
