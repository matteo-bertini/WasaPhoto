package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/utils"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Controllo che sia stato specificato l'authcode
	auth_header := r.Header.Get("Authorization")
	if auth_header == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		// Authcode specificato
		token := utils.ParseAuthToken(auth_header)
		if token == nil {
			// è stata passata una stringa vuota come authentication
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			// L'authcode è stato passato,va controllato che sia corretto per autorizzare l'operazione
			var to_del string
			to_del = strings.Split(r.URL.Path, "/")[4]

			// Controllo dell'authcode
			authorized, err := rt.db.Authcheck(to_del, *token)
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
					// L'user è autorizzato a fare l'operazione

					username := strings.Split(r.URL.Path, "/")[2]
					err = rt.db.UnfollowUser(username, to_del)
					if err != nil {
						if err.Error() == "The username searched does not exists." {
							w.WriteHeader(http.StatusNotFound)
							return
						} else {
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
					} else {
						w.WriteHeader(http.StatusNoContent)
						return
					}

				}
			}

		}
	}
}
