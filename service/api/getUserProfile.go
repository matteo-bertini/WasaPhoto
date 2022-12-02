package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/database"
	"WasaPhoto/service/utils"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user *database.User
	auth_header := r.Header.Get("Authorization")
	// Authentication token not specified in the header,sending back BadRequest
	if auth_header == "" {
		w.WriteHeader(http.StatusBadRequest)
		return

	} else { // Authentication token specified in the header.

		// Authentication token extraction
		token := utils.ParseAuthToken(auth_header)
		if token == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} else {
			if r.URL.Query().Has("username") {
				var username string = r.URL.Query().Get("username")
				var err error
				check := utils.CheckUsername(username)
				if check == false {
					w.WriteHeader(http.StatusBadRequest)
					return
				} else {
					// Authorization Token Check
					// in questo caso devo controllare solo che colui che fa la richiesta sia iscritto,quindi che il token che mi passa
					// esista nel database
					var auth_res *bool
					auth_res, err = rt.db.GetUserProfile_Authcheck(*token)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					if *auth_res == false {
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
					user, err = rt.db.GetUserProfile(username)
					if err != nil {
						// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
						// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
						ctx.Logger.WithError(err).Error("Can't find the user searched.")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(*user)

				}
			} else {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

		}
	}

}
