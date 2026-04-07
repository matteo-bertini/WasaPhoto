package api

import (
	"WasaPhoto/service/api/reqcontext"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// AuthMiddleware validates the token and adds the UserID to the context
func (rt *_router) AuthMiddleware(fn httpRouterHandler) httpRouterHandler {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
		// 1. Get the token from the header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.Logger.Error("No valid Bearer token provided")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// 2. Ask the database who owns this token
		userID, err := rt.db.GetUserIDByToken(token)
		if err != nil {
			ctx.Logger.WithError(err).Error("Invalid or expired token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// 3. Save the ID in our custom context
		ctx.UserID = userID

		// 4. Continue the chain
		fn(w, r, ps, ctx)
	}
}
