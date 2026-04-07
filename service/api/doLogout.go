package api

import (
	"WasaPhoto/service/api/reqcontext"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// doLogout handles the session removal
func (rt *_router) doLogout(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// The UserID is already populated by the AuthMiddleware
	userID := ctx.UserID

	// Perform the logout in the database
	err := rt.db.DoLogout(userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("Database error during logout")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Success response
	w.WriteHeader(http.StatusNoContent)
}
