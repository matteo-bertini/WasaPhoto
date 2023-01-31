package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/utils"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// First step : checking if the username in the request body and the authcode passed are consistent
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
			// Authcode passed,we extract the username from request body for the check
			var req_body followUserRequestBody
			err := json.NewDecoder(r.Body).Decode(&req_body)
			if err != nil {
				// Il body passato non è parsabile come JSON,lo rifiuto e mando BadRequest
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			FollowerId := req_body.FollowerId
			res, err := rt.db.FollowUser_Authcheck(FollowerId, *token)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			} else {
				if *res == false {
					w.WriteHeader(http.StatusUnauthorized)
					return
				} else {
					// L'authcode passato e l'username passato sono coerenti,ora va aggiunto il follower
					username := strings.Split(r.URL.Path, "/")[2]
					err := rt.db.FollowUser(username, FollowerId)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						return
					} else {
						w.WriteHeader(http.StatusCreated)
						w.Header().Set("Content-Type", "application/json")
						_ = json.NewEncoder(w).Encode(req_body)
						return

					}

				}
			}
		}
	}

}
