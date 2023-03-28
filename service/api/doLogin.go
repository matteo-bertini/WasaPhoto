package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"

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

		// C'è un errore nel RequestBody passato (nomi dei campi errati,campi non specificati,ecc)
		if len(doLoginRequestBody.Username) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			ctx.Logger.Error("Il json nel RequestBody presenta degli errori.")
			return

		} else { // Il RequestBody passato non presenta errori

			// Controllo che l'Username passato nel RequestBody sia una stringa conforme alle specifiche
			if !utils.CheckUsername(doLoginRequestBody.Username) {
				w.WriteHeader(http.StatusBadRequest)
				ctx.Logger.Error("L'Username passato nel RequestBody non è conforme alle specifiche.")
				return
			} else {
				// L'Username passato nel RequestBody è conforme alle specifiche progettuali
				id, err := rt.db.DoLogin(doLoginRequestBody.Username)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					ctx.Logger.WithError(err).Error("Si è verificato un errore nelle operazioni sul database.")
					return
				} else {
					// Imposto l'header della risposta e scrivo lo status 201 Created
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusCreated)

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
