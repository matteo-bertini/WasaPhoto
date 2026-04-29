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
	rt.router.DELETE("/users/:username", rt.wrap(rt.AuthMiddleware(rt.DeleteUserHandler)))
	rt.router.PUT("/users/:username", rt.wrap(rt.AuthMiddleware(rt.UpdateUsernameHandler)))

	// Social relationships management
	rt.router.GET("/users/:username/followers", rt.wrap(rt.AuthMiddleware(rt.GetFollowersListHandler)))
	rt.router.GET("/users/:username/following", rt.wrap(rt.AuthMiddleware(rt.GetFollowingListHandler)))
	rt.router.PUT("/users/:username/followers", rt.wrap(rt.AuthMiddleware(rt.FollowUserHandler)))
	rt.router.DELETE("/users/:username/followers", rt.wrap(rt.AuthMiddleware(rt.UnfollowUserHandler)))
	rt.router.GET("/users/:username/ban", rt.wrap(rt.AuthMiddleware(rt.GetBanListHandler)))
	rt.router.PUT("/users/:username/ban", rt.wrap(rt.AuthMiddleware(rt.BanUserHandler)))
	rt.router.DELETE("/users/:username/ban", rt.wrap(rt.AuthMiddleware(rt.UnbanUserHandler)))

	// Image uploading and posts management
	rt.router.POST("/users/:username/posts", rt.wrap(rt.AuthMiddleware(rt.UploadPostHandler)))
	rt.router.GET("/users/:username/posts/:postId", rt.wrap(rt.GetPhotoHandler))
	rt.router.DELETE("/users/:username/posts/:postId", rt.wrap(rt.AuthMiddleware(rt.DeletePostHandler)))
	rt.router.PUT("/users/:username/posts/:postId/likes", rt.wrap(rt.AuthMiddleware(rt.LikePostHandler)))
	rt.router.DELETE("/users/:username/posts/:postId/likes", rt.wrap(rt.AuthMiddleware(rt.UnlikePostHandler)))
	rt.router.GET("/users/:username/posts/:postId/likes", rt.wrap(rt.AuthMiddleware(rt.GetLikesHandler)))

	// Comment Routes
	// Registra un nuovo commento: POST /users/:username/posts/:postId/comments
	rt.router.POST("/users/:username/posts/:postId/comments", rt.wrap(rt.AuthMiddleware(rt.AddCommentHandler)))

	// Recupera la lista dei commenti: GET /users/:username/posts/:postId/comments
	rt.router.GET("/users/:username/posts/:postId/comments", rt.wrap(rt.AuthMiddleware(rt.GetCommentsHandler)))

	// Elimina un commento specifico: DELETE /users/:username/posts/:postId/comments/:commentId
	rt.router.DELETE("/users/:username/posts/:postId/comments/:commentId", rt.wrap(rt.AuthMiddleware(rt.DeleteCommentHandler)))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
