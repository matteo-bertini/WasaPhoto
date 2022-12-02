package api

import (
	"WasaPhoto/service/api/reqcontext"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) addUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Read the new content for the User to add from the request body.
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		// The body was not a parseable JSON, reject it
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user passed and parsed is valid.
	if user.IsValid() == false {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Create the user in the database. Note that this function will return a new instance of the user.
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
