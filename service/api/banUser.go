package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/utils"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Estrazione del token di autenticazione dall'header
	auth_header := r.Header.Get("Authorization")
	if auth_header == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		// Here we parse the authcode from the header
		token := utils.ParseAuthToken(auth_header)
		if token == nil {
			// This is the case in which we have passed an empty string as authcode
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			// Il token è stato passato e va controllato che sia consistente
			username := strings.Split(r.URL.Path, "/")[2]
			authorized, err := rt.db.Authcheck(username, *token)
			if err != nil {
				if err.Error() == "The username searched does not exists." {
					w.WriteHeader(http.StatusNotFound)
					return

				} else {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			} else {
				if *authorized == false {
					w.WriteHeader(http.StatusUnauthorized)
					return
				} else {
					//TODO
					w.WriteHeader(http.StatusOK)
					return

				}
			}
		}
	}

}
