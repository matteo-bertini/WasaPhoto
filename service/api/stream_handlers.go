package api

import (
	"WasaPhoto/service/api/reqcontext"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) GetStreamHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Extract and parse pagination parameters
	rawLimit := r.URL.Query().Get("limit")
	rawOffset := r.URL.Query().Get("offset")

	const maxLimit = 100

	limit, err := strconv.Atoi(rawLimit)
	if err != nil || limit <= 0 {
		limit = 20 // Default limit if not specified or invalid
	}
	if limit > maxLimit {
		limit = maxLimit // Cap to avoid unbounded result sets
	}

	offset, err := strconv.Atoi(rawOffset)
	if err != nil || offset < 0 {
		offset = 0 // Default offset if not specified or invalid
	}

	// 2. Call the database function
	// We pass the requester's ID (from context), limit, and offset
	stream, err := rt.db.GetStream(ctx.UserID, limit, offset)
	if err != nil {
		ctx.Logger.WithError(err).Error("GetStreamHandler: failed to fetch hybrid stream from database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stream); err != nil {
		ctx.Logger.WithError(err).Error("GetStreamHandler: failed to encode response JSON")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
