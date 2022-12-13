package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/database"
	"WasaPhoto/service/utils"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) deleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Parsing the username in the path
	username := ps.ByName("Username")
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
			if utils.CheckUsername(username) == false {
				w.WriteHeader(http.StatusBadRequest)
				return

			}
			res, err := rt.db.DeleteUser_Authcheck(username, *token)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				if *res == false {
					w.WriteHeader(http.StatusUnauthorized)
					return
				} else {
					err = rt.db.DeleteUser(username)
					if errors.Is(err, database.ErrUserDoesNotExist) {
						// The fountain (indicated by `id`) does not exist, reject the action indicating an error on the client side.
						w.WriteHeader(http.StatusNotFound)
						return
					} else if err != nil {
						// In this case, we have an error on our side. Log the error (so we can be notified) and send a 500 to the user
						// Note: we are using the "logger" inside the "ctx" (context) because the scope of this issue is the request.
						// Note (2): we are adding the error and an additional field (`id`) to the log entry, so that we will receive
						// the identifier of the fountain that triggered the error.
						ctx.Logger.WithError(err).WithField("username", username).Error("can't delete the user")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}

				}
			}
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
