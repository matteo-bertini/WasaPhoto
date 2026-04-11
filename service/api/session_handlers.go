package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/models"

	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// LoginHandler handles the session start
func (rt *_router) LoginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// 1) Decode the request body.
	var req models.DoLoginRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// Handle empty body or decoding errors.
		if errors.Is(err, io.EOF) {
			ctx.Logger.WithError(err).Error("LoginHandler: Missing request body")
		} else {
			ctx.Logger.WithError(err).Error("LoginHandler: Failed to decode request body")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 2) Validate required fields.
	if req.Username == "" || req.Password == "" || req.IsSignUp == nil {
		ctx.Logger.WithError(models.ErrMissingFields).Error("LoginHandler: request body validation failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 3) Check if username and password strings are valid according to specs.
	if !models.CheckUsername(req.Username) || !models.CheckPassword(req.Password) {
		ctx.Logger.WithError(models.ErrInvalidUsernameOrPassword).Error("LoginHandler: username/password validation failed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 4) Execute the unified database logic.
	sessionToken, err := rt.db.DoLogin(req.Username, req.Password, *req.IsSignUp)
	if err != nil {
		// Handle specific authentication errors.
		if errors.Is(err, models.ErrInvalidCredentials) {
			ctx.Logger.WithError(err).Error("LoginHandler: login process failed ")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if errors.Is(err, models.ErrUserAlreadyExists) {
			ctx.Logger.WithError(err).Error("LoginHandler: login process failed")
			w.WriteHeader(http.StatusConflict)
			return
		}

		// Handle generic database errors.
		ctx.Logger.WithError(err).Error("LoginHandler: login process failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 5) Set response headers and status code.
	w.Header().Set("Content-Type", "application/json")
	if *req.IsSignUp {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	// 6) Encode and send the response body.
	var res models.DoLoginResponseBody
	res.SessionToken = *sessionToken
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		ctx.Logger.WithError(err).Error("LoginHandler: failed to encode response JSON")
	}

}

// LogoutHandler handles the session removal
func (rt *_router) LogoutHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// The UserID is already populated by the AuthMiddleware
	userID := ctx.UserID

	// Perform the logout in the database
	err := rt.db.DoLogout(userID)
	if err != nil {
		ctx.Logger.WithError(err).Error("LogoutHandler: logout process failed")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Success response
	w.WriteHeader(http.StatusNoContent)
}
