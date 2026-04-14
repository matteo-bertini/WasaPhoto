package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/models"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// FollowUserHandler handles the PUT request to follow a user
func (rt *_router) FollowUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	targetUsername := ps.ByName("username")
	actorUsername := ps.ByName("actor_username")

	// 1. Security Cross-Check
	// We compare the ID from the context (auth token) with the ID of the username in the URL.
	// This prevents a logged-in user from performing actions on behalf of others.

	authUserID := ctx.UserID // From Authentication Middleware

	actorID, err := rt.db.GetIDByUsername(actorUsername)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			ctx.Logger.WithError(err).Error("FollowUserHandler: failed to resolve actor ID")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("FollowUserHandler: failed to resolve actor ID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if actorID != authUserID {
		ctx.Logger.Warn("FollowUserHandler: security breach attempt - userID mismatch")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 2. Resolve Target Username to ID
	targetID, err := rt.db.GetIDByUsername(targetUsername)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			ctx.Logger.WithError(err).Error("FollowUserHandler: failed to resolve target ID")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("FollowUserHandler: failed to resolve target ID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. Database Action
	// The DB layer handles business logic (e.g., preventing self-follows or checking bans)
	err = rt.db.FollowUser(actorID, targetID)
	if err != nil {
		if errors.Is(err, models.ErrForbiddenAction) || errors.Is(err, models.ErrSelfFollow) {
			ctx.Logger.WithError(err).Error("FollowUserHandler: failed to perform follow action")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		ctx.Logger.WithError(err).Error("FollowUserHandler: failed to perform follow action")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 4. Success Response
	w.WriteHeader(http.StatusNoContent)

}

// UnfollowUserHandler handles the DELETE request to unfollow a user
func (rt *_router) UnfollowUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	targetUsername := ps.ByName("username")
	actorUsername := ps.ByName("actor_username")

	// 1. Security Cross-Check
	// We compare the ID from the context (auth token) with the ID of the username in the URL.
	// This prevents a logged-in user from performing actions on behalf of others.

	authUserID := ctx.UserID // From Authentication Middleware

	actorID, err := rt.db.GetIDByUsername(actorUsername)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			ctx.Logger.WithError(err).Error("UnfollowUserHandler: failed to resolve actor ID")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("UnfollowUserHandler: failed to resolve actor ID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if actorID != authUserID {
		ctx.Logger.Warn("UnfollowUserHandler: security breach attempt - userID mismatch")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 2. Resolve Target Username to ID
	targetID, err := rt.db.GetIDByUsername(targetUsername)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			ctx.Logger.WithError(err).Error("UnfollowUserHandler: failed to resolve target ID")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("UnfollowUserHandler: failed to resolve target ID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = rt.db.UnfollowUser(actorID, targetID)
	if err != nil {
		ctx.Logger.WithError(err).Error("UnfollowUserHandler: failed to perform unfollow action")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// BanUserHandler handles the PUT request to ban a user
func (rt *_router) BanUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	actorUsername := ps.ByName("username")
	targetUsername := ps.ByName("target_username")

	// 1. Security check: Identity verification
	actorID, err := rt.db.GetIDByUsername(actorUsername)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			ctx.Logger.WithError(err).Error("BanUserHandler: failed to resolve actor ID")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("BanUserHandler: failed to resolve actor ID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if actorID != ctx.UserID {
		ctx.Logger.Warn("BanUserHandler: unauthorized ban attempt")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 2. Resolve target ID
	targetID, err := rt.db.GetIDByUsername(targetUsername)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			ctx.Logger.WithError(err).Error("BanUserHandler: failed to resolve target ID")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("BanUserHandler: failed to resolve target ID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. Database Action
	err = rt.db.BanUser(actorID, targetID)
	if err != nil {
		if errors.Is(err, models.ErrSelfBan) {
			ctx.Logger.WithError(err).Error("BanUserHandler: failed to perform ban action")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		ctx.Logger.WithError(err).Error("BanUserHandler: database error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UnbanUserHandler handles the DELETE request to unban a user
func (rt *_router) UnbanUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	actorUsername := ps.ByName("username")
	targetUsername := ps.ByName("target_username")

	// 1. Security Cross-Check
	// We compare the ID from the context (auth token) with the ID of the username in the URL.
	// This prevents a logged-in user from performing actions on behalf of others.

	authUserID := ctx.UserID // From Authentication Middleware

	actorID, err := rt.db.GetIDByUsername(actorUsername)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			ctx.Logger.WithError(err).Error("UnbanUserHandler: failed to resolve actor ID")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("UnbanUserHandler: failed to resolve actor ID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if actorID != authUserID {
		ctx.Logger.Warn("UnbanUserHandler: security breach attempt - userID mismatch")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 2. Resolve Target Username to ID
	targetID, err := rt.db.GetIDByUsername(targetUsername)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			ctx.Logger.WithError(err).Error("UnbanUserHandler: failed to resolve target ID")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("UnbanUserHandler: failed to resolve target ID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. Database Action (Idempotent DELETE)
	err = rt.db.UnbanUser(actorID, targetID)
	if err != nil {
		ctx.Logger.WithError(err).Error("UnbanUserHandler: database error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
