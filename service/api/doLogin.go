package api

import (
	"WasaPhoto/service/api/reqcontext"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req_body dologinRequestBody
	var ret *string
	err := json.NewDecoder(r.Body).Decode(&req_body)
	if err != nil {
		// Il body passato non è parsabile come JSON,lo rifiuto e mando BadRequest
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		// Ricerca nel database dell'authstring per l'username specificato (se non esiste verrà creata)
		ret, err = rt.db.DoLogin(req_body.Username)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return

		} else {
			w.WriteHeader(http.StatusCreated)
			resp_body := doLoginResponseBody{Identifier: *ret}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp_body)
			return

		}

	}

}
