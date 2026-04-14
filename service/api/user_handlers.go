package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/models"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) GetUserProfileHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// 1. Extract the Target Username from the URL path
	targetUsername := ps.ByName("username")

	// 2. Retrieve the Requester's ID from the context (populated by the Auth Middleware)
	// We assume the middleware stores the ID as a string under the key "userID"
	requestingUserID := ctx.UserID

	// 3. Database Layer Call
	// The DB method handles banning logic, relationship checks (IsFollowing),
	// and gathers profile stats + post list in an optimized way.
	profile, err := rt.db.GetUserProfile(targetUsername, requestingUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// User does not exist
			ctx.Logger.WithError(models.ErrUserNotFound).Error("GetUserProfileHandler: failed to fetch profile")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, models.ErrProfileAccessForbidden) {
			// Handle banning relationship between the two users
			ctx.Logger.WithError(models.ErrProfileAccessForbidden).Error("GetUserProfileHandler: failed to fetch profile")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		// Generic server or database error
		rt.baseLogger.WithError(err).Error("GetUserProfileHandler: failed to fetch profile")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 4. Send the JSON response to the Frontend
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(profile)
}
