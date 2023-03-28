package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/segmentio/ksuid"
)

func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Controllo che l'username nell'URL corrisponda ad un user esistente
	urlusername := strings.Split(r.URL.Path, "/")[2]
	err := rt.db.CheckUserExistence(urlusername)
	if err != nil {
		if errors.Is(err, utils.ErrUserDoesNotExist) {
			w.WriteHeader(http.StatusNotFound)
			ctx.Logger.WithError(err).Error("L'user specificato nell'URL non esiste.")
			return

		} else {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("Si è verificato un errore nel controllare l'esistenza dell'user nell'URL.")
			return

		}
	} else {
		// L'username specificato nell'URL corrisponde ad un user esistente
		err = rt.db.CheckAuthorization(r, urlusername)
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
			// L'id passato nell'Authorization è quello dell'user nell'URL

			// Controllo che il RequestBody sia stato specificato correttamente
			if r.ContentLength == 0 || r.Header.Get("Content-Type") != "image/png" {
				w.WriteHeader(http.StatusBadRequest)
				ctx.Logger.Error("La richiesta non ha specificato correttamente il RequestBody.")
				return
			} else {
				// RequestBody passato correttamente
				buf, err := io.ReadAll(r.Body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					ctx.Logger.WithError(err).Error("Si è verificato un errore nel parsing della foto nel RequestBody.")
					return
				} else {
					var osfile *os.File

					// Verrà creata una cartella con nome Id dell'user se già non esiste (per questo non viene controllato l'errore)
					path := "/tmp/WasaPhoto/" + strings.Split(r.Header.Get("Authorization"), " ")[1]
					_ = os.Mkdir(path, os.ModePerm)

					// Creazione dell'Id della foto
					photoid := ksuid.New().String()

					// Creazione del file
					osfile, err = os.Create(path + "/" + photoid + ".png")
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						ctx.Logger.WithError(err).Error("Si è verificato un errore nella creazione del file .png da salvare.")
						return
					} else {
						// Scrittura dei bytes nel file appena creato
						_, err = osfile.Write(buf)
						if err != nil {
							w.WriteHeader(http.StatusInternalServerError)
							ctx.Logger.WithError(err).Error("Si è verificato un errore nella scrittura sul file .png .")
							return
						} else {

							UploadPhotoResponseBody := Photo{PhotoId: photoid, LikesNumber: 0, CommentsNumber: 0, DateOfUpload: time.Now().String()}
							err = rt.db.UploadPhoto(UploadPhotoResponseBody.PhotoToDatabase(), strings.Split(r.Header.Get("Authorization"), " ")[1])
							if err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								ctx.Logger.WithError(err).Error("Si è verificato un errore nelle operazioni di database.")
								return
							} else {
								w.Header().Set("Content-Type", "application/json")
								w.WriteHeader(http.StatusCreated)
								err = json.NewEncoder(w).Encode(UploadPhotoResponseBody)
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
