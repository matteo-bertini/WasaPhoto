package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/utils"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) addUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	auth_header := r.Header.Get("Authorization")
	// Authentication token not specified in the header,sending back BadRequest
	if auth_header == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else {
		// Extraction of authentication token specified in the header
		token := utils.ParseAuthToken(auth_header)
		if token == nil {
			w.WriteHeader(http.StatusBadRequest)

		} else {
			var req_body addUserRequestBody
			err := json.NewDecoder(r.Body).Decode(&req_body)
			if err != nil {
				// Il body passato non è parsabile come JSON,lo rifiuto e mando BadRequest
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			username := req_body.Username
			if utils.CheckUsername(username) == false {
				w.WriteHeader(http.StatusBadRequest)
				return

			}
			res, err := rt.db.AddUser_Authcheck(username, *token)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				if *res == false {
					w.WriteHeader(http.StatusUnauthorized)
					return
				} else {
					var user User
					user.Username = username
					dbuser, err := rt.db.AddUser(user.ToDatabase())
					if err != nil {
						// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
						// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
						ctx.Logger.WithError(err).Error("Can't create the user")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

					// Here we can re-use `user` as FromDatabase is overwriting every variabile in the structure.
					user.FromDatabase(dbuser)

					// Send the output to the user.
					w.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(w).Encode(user)
				}

			}
		}

	}

}
