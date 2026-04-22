package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/models"
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
		// This aligns the backend with the OpenAPI specification
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
			ctx.Logger.WithError(err).Error("AuthMiddleware: database could not validate token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// 5. Inject the validated UserID into the request context for downstream handlers
		ctx.UserID = userID

		// 6. Proceed with the request chain
		fn(w, r, ps, ctx)
	}
}
