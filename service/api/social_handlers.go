package api

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// FollowUserHandler handles the PUT request to follow a user
func (rt *_router) FollowUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	targetUsername := ps.ByName("target_username")
	actorUsername := ps.ByName("actor_username")

	// 1. Security Check: actor_username in path must match the Authenticated User
	// We assume the middleware stores the Username or we fetch it via ID
	authUserID := r.Context().Value("userID").(string)

	// Verification logic: Fetch Actor ID from Username to ensure they match
	actorID, err := rt.db.GetIDByUsername(actorUsername)
	if err != nil || actorID != authUserID {
		http.Error(w, "Unauthorized: identity mismatch", http.StatusUnauthorized)
		return
	}

	// 2. Resolve Target Username to ID
	targetID, err := rt.db.GetIDByUsername(targetUsername)
	if err != nil {
		http.Error(w, "Target user not found", http.StatusNotFound)
		return
	}

	// 3. Database Action
	err = rt.db.FollowUser(actorID, targetID)
	if err != nil {
		// Distinguish between business logic error (Self-follow/Ban) and DB error
		if strings.Contains(err.Error(), "forbidden") || strings.Contains(err.Error(), "yourself") {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 4. Success: 204 No Content
	w.WriteHeader(http.StatusNoContent)
}

// UnfollowUserHandler handles the DELETE request to unfollow a user
func (rt *Router) UnfollowUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	targetUsername := ps.ByName("target_username")
	actorUsername := ps.ByName("actor_username")

	authUserID := r.Context().Value("userID").(string)
	actorID, _ := rt.db.GetIDByUsername(actorUsername)

	if actorID != authUserID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	targetID, err := rt.db.GetIDByUsername(targetUsername)
	if err != nil {
		http.Error(w, "Target user not found", http.StatusNotFound)
		return
	}

	err = rt.db.UnfollowUser(actorID, targetID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
