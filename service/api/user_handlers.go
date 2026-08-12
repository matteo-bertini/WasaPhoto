package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/models"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// GetUserProfileHandler processes the request to retrieve a user's full profile.
// It extracts the target username from the URI and the requester's identity from the request context.
func (rt *_router) GetUserProfileHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// 1. Extract the Target Username from the URL path
	targetUsername := ps.ByName("username")

	// 2. Retrieve the Requester's ID from the context (populated by the Auth Middleware)
	requestingUserID := ctx.UserID

	// 3. Database Layer Call
	// The DB method handles banning logic, relationship checks (IsFollowing),
	// and gathers profile stats + post list in an optimized way.
	profile, err := rt.db.GetUserProfile(targetUsername, requestingUserID)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			// User does not exist
			ctx.Logger.WithError(err).Error("GetUserProfileHandler: failed to fetch profile")
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

// UpdateUsernameHandler changes the username of the authenticated user.
func (rt *_router) UpdateUsernameHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Extract the current username from the URI path
	pathUsername := ps.ByName("username")

	// 2. Authorization: Get the ID of the user in the path
	targetUserID, err := rt.db.GetIDByUsername(pathUsername)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("UpdateUsername: database error during identity check")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. Security Check: Only the account owner can change their own username
	if targetUserID != ctx.UserID {
		ctx.Logger.Warn("UpdateUsername: forbidden attempt to change another user's name")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 4. Parse the request body
	var body struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		ctx.Logger.WithError(err).Warn("UpdateUsername: invalid JSON body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate against the same pattern used at signup (minLength, maxLength, charset)
	newUsername := strings.TrimSpace(body.Username)
	if !models.CheckUsername(newUsername) {
		ctx.Logger.Warn("UpdateUsername: username format validation failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 5. Update in Database
	err = rt.db.UpdateUsername(ctx.UserID, newUsername)
	if err != nil {
		if errors.Is(err, models.ErrUsernameTaken) {
			// YAML 409 Conflict
			w.WriteHeader(http.StatusConflict)
			return
		}
		ctx.Logger.WithError(err).Error("UpdateUsername: DB update failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 6. Success (204 No Content)
	ctx.Logger.Infof("User %s successfully changed to %s", pathUsername, newUsername)
	w.WriteHeader(http.StatusNoContent)
}

// DeleteUserHandler handles the permanent deletion of a user account and all related data.
// It removes the database record (triggering CASCADE) and deletes the user's media folder.
func (rt *_router) DeleteUserHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Extract the username from the URL path
	username := ps.ByName("username")

	// 2. Security Check: Retrieve the target user's ID by their username
	targetUserID, err := rt.db.GetIDByUsername(username)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			ctx.Logger.WithField("username", username).Warn("DeleteUserHandler: user not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Generic database error (e.g. connection issue): must not fall through
		// to the authorization check below, otherwise a DB failure would be
		// misreported as a 403 Forbidden instead of a 500 Internal Server Error.
		ctx.Logger.WithError(err).Error("DeleteUserHandler: database error during identity check")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. Authorization: Ensure the requester is only deleting their own account
	// We compare the ID from the URL with the ID extracted from the Auth Token
	if targetUserID != ctx.UserID {
		ctx.Logger.WithFields(logrus.Fields{
			"requester": ctx.UserID,
			"target":    targetUserID,
		}).Warn("DeleteUserHandler: unauthorized deletion attempt")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 4. Database Deletion
	// This will trigger ON DELETE CASCADE on: profiles, posts, comments, likes, follows, bans
	err = rt.db.DeleteUser(ctx.UserID)
	if errors.Is(err, models.ErrUserNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		ctx.Logger.WithError(err).Error("DeleteUserHandler: failed to delete user from database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 5. Filesystem Cleanup
	// Since we organized posts in /UploadDir/userID/, we can delete the whole folder at once
	userMediaDir := filepath.Join(UploadDir, ctx.UserID)

	// os.RemoveAll is the best choice here: it deletes the directory and all its contents
	err = os.RemoveAll(userMediaDir)
	if err != nil {
		// We log the error but don't return 500 because the account is already gone from the DB.
		// Failing to delete a file is a maintenance issue, not a request failure.
		ctx.Logger.WithError(err).WithField("path", userMediaDir).Warn("DeleteUserHandler: could not remove media directory")
	}

	// 6. Final Response
	ctx.Logger.WithField("userID", ctx.UserID).Info("DeleteUserHandler: account and data successfully deleted")
	w.WriteHeader(http.StatusNoContent)
}
