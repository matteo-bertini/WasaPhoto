package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/utils"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// 1) Decode the request body.
	var req doLoginRequestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// Handle empty body or decoding errors.
		if errors.Is(err, io.EOF) {
			ctx.Logger.WithError(err).Error("Missing request body.")
		} else {
			ctx.Logger.WithError(err).Error("Failed to decode request body.")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 2) Validate required fields.
	if req.Username == "" || req.Password == "" || req.IsSignUp == nil {
		ctx.Logger.Error("Invalid or missing fields in the JSON body.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 3) Check if username and password strings are valid according to specs.
	if !utils.CheckUsername(req.Username) || !utils.CheckPassword(req.Password) {
		ctx.Logger.Error("Username or Password does not meet the requirements.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 4) Execute the unified database logic.
	id, err := rt.db.DoLogin(req.Username, req.Password, *req.IsSignUp)
	if err != nil {
		// Handle specific authentication errors.
		if errors.Is(err, utils.ErrInvalidCredentials) {
			ctx.Logger.WithError(err).Error("Invalid credentials provided.")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if errors.Is(err, utils.ErrUserAlreadyExists) {
			ctx.Logger.WithError(err).Error("User already exists.")
			w.WriteHeader(http.StatusConflict)
			return
		}

		// Handle generic database errors.
		ctx.Logger.WithError(err).Error("Database operation failed.")
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
	var response doLoginResponseBody
	response.Identifier = *id
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		ctx.Logger.WithError(err).Error("Failed to encode response JSON.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
