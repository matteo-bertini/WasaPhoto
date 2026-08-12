package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/models"
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware validates the Bearer token and injects the UserID into the request context.
func (rt *_router) AuthMiddleware(fn httpRouterHandler) httpRouterHandler {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

		// 1. Retrieve the Authorization header
		authHeader := r.Header.Get("Authorization")

		// 2. Formal check: the header must exist and start with the "Bearer " prefix
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.Logger.WithFields(logrus.Fields{
				"header": authHeader,
				"error":  models.ErrUnauthorized.Error(),
			}).Error("AuthMiddleware: missing or invalid Bearer prefix")

			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// 3. Extract the token by removing the "Bearer " prefix (the first 7 characters)
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Ensure the token is not empty after removing the prefix
		if token == "" {
			ctx.Logger.WithError(models.ErrUnauthorized).Error("AuthMiddleware: token is empty after Bearer prefix")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// 4. Validate the token against the database
		// GetUserIDByToken verifies if the token is valid and returns the associated UserID
		userID, err := rt.db.GetUserIDByToken(token)
		if err != nil {
			if errors.Is(err, models.ErrInvalidToken) {
				// The token is genuinely missing/invalid: this is a client-side 401.
				ctx.Logger.WithError(err).Error("AuthMiddleware: invalid or expired token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// Any other error (e.g. a DB/connection failure) is an infrastructure
			// problem, not an authorization decision: report it as a 500 so it
			// isn't masked as "unauthorized".
			ctx.Logger.WithError(err).Error("AuthMiddleware: database error while validating token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// 5. Inject the validated UserID into the request context for downstream handlers
		ctx.UserID = userID

		// 6. Proceed with the request chain
		fn(w, r, ps, ctx)
	}
}
