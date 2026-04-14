package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Account and session management //
	rt.router.POST("/session", rt.wrap(rt.LoginHandler))
	rt.router.DELETE("/session", rt.wrap(rt.AuthMiddleware(rt.LogoutHandler)))

	// User management
	rt.router.GET("/users/:username", rt.wrap(rt.AuthMiddleware(rt.GetUserProfileHandler)))

	// Social relationships management
	rt.router.PUT("/users/:username/followers/:actor_username", rt.wrap(rt.AuthMiddleware(rt.FollowUserHandler)))
	rt.router.DELETE("/users/:username/followers/:actor_username", rt.wrap(rt.AuthMiddleware(rt.UnfollowUserHandler)))
	rt.router.PUT("/users/:username/bans/:target_username", rt.wrap(rt.AuthMiddleware(rt.BanUserHandler)))
	rt.router.DELETE("/users/:username/bans/:target_username", rt.wrap(rt.AuthMiddleware(rt.UnbanUserHandler)))

	// deleteUser //
	//rt.router.DELETE("/users/:Username/", rt.wrap(rt.deleteUser))

	// getMyStream //
	//rt.router.GET("/users/:Username/", rt.wrap(rt.getMyStream))

	// setMyUsername //
	//rt.router.PUT("/users/:Username/username", rt.wrap(rt.setMyUsername))

	// getBanned //
	//rt.router.GET("/users/:Username/bannedusers/", rt.wrap(rt.getBanned))

	// uploadPhoto //
	//rt.router.POST("/users/:Username/photos/", rt.wrap(rt.uploadPhoto))

	// getPhoto //
	/*rt.router.GET("/users/:Username/photos/:PhotoId/", rt.wrap(rt.getPhoto))

	// deletePhoto //
	rt.router.DELETE("/users/:Username/photos/:PhotoId/", rt.wrap(rt.deletePhoto))

	// getLikes //
	rt.router.GET("/users/:Username/photos/:PhotoId/likes/", rt.wrap(rt.getLikes))

	// likePhoto //
	rt.router.POST("/users/:Username/photos/:PhotoId/likes/", rt.wrap(rt.likePhoto))

	// unlikePhoto //
	rt.router.DELETE("/users/:Username/photos/:PhotoId/likes/:LikeId", rt.wrap(rt.unlikePhoto))

	// getComments //
	rt.router.GET("/users/:Username/photos/:PhotoId/comments/", rt.wrap(rt.getComments))

	// commentPhoto //
	rt.router.POST("/users/:Username/photos/:PhotoId/comments/", rt.wrap(rt.commentPhoto))

	// uncommentPhoto //
	rt.router.DELETE("/users/:Username/photos/:PhotoId/comments/:CommentId", rt.wrap(rt.uncommentPhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness) */

	return rt.router
}
