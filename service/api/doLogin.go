package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Controllo che la richiesta abbia specificato RequestBody ed in tal caso lo estraggo

	var doLoginRequestBody doLoginRequestBody
	err := json.NewDecoder(r.Body).Decode(&doLoginRequestBody)
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

		// C'è un errore nel RequestBody passato (nomi dei campi errati,campi necessari non specificati,ecc)
		if len(doLoginRequestBody.Username) == 0 || len(doLoginRequestBody.Password) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Error("Il json nel RequestBody presenta degli errori.")
			return

		} else { // Il RequestBody passato non presenta errori

			// Controllo che l'Username e la Password passati nel RequestBody sia una stringa conforme alle specifiche
			if !(utils.CheckUsername(doLoginRequestBody.Username) && utils.CheckPassword(doLoginRequestBody.Password)) {
				w.WriteHeader(http.StatusBadRequest)
				ctx.Logger.Error("L'Username e/o la Password passati nel RequestBody non sono conformi alle specifiche.")
				return
			} else {
				// L'Username e la Password passati nel RequestBody sono conformi alle specifiche progettuali
				id, err, created := rt.db.DoLogin(doLoginRequestBody.Username, doLoginRequestBody.Password)
				if err != nil {
					if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
						w.WriteHeader(http.StatusUnauthorized)
						ctx.Logger.WithError(err).Error("Credenziali non valide.")
						return

					} else {
						w.WriteHeader(http.StatusInternalServerError)
						ctx.Logger.WithError(err).Error("Si è verificato un errore nelle operazioni sul database.")
						return
					}
				} else {
					// Imposto l'header della risposta e scrivo lo status 201 Created se è stato creato un nuovo utente o 200 OK altrimenti
					w.Header().Set("Content-Type", "application/json")
					if *created == true {
						w.WriteHeader(http.StatusCreated)

					} else {
						w.WriteHeader(http.StatusOK)

					}

					// Faccio l'encoding del ResponseBody per mandarlo nel json di risposta
					var doLoginResponseBody doLoginResponseBody
					doLoginResponseBody.Identifier = *id
					err = json.NewEncoder(w).Encode(doLoginResponseBody)

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
