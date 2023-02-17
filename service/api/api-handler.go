package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes

	// doLogin //
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// addUser //
	rt.router.POST("/users/", rt.wrap(rt.addUser))

	// getUserProfile
	rt.router.GET("/users/", rt.wrap(rt.getUserProfile))

	// deleteUser //
	rt.router.DELETE("/users/:Username/", rt.wrap(rt.deleteUser))

	rt.router.GET("/users/:Username/", rt.wrap(rt.getMyStream))

	// setMyUsername //
	rt.router.PUT("/users/:Username/username", rt.wrap(rt.setMyUsername))

	// followUser //
	rt.router.POST("/users/:Username/followers/", rt.wrap(rt.followUser))

	// unfollowUser //
	rt.router.DELETE("/users/:Username/followers/:FollowerId", rt.wrap(rt.unfollowUser))

	// banUser //
	rt.router.POST("/users/:Username/bannedusers/", rt.wrap(rt.banUser))

	// unbanUser //
	rt.router.DELETE("/users/:Username/bannedusers/:BannedId", rt.wrap(rt.unbanUser))

	// uploadPhoto //
	rt.router.POST("/users/:Username/photos/", rt.wrap(rt.uploadPhoto))

	// getPhoto //
	rt.router.GET("/users/:Username/photos/:PhotoId/", rt.wrap(rt.getPhoto))

	// deletePhoto //
	rt.router.DELETE("/users/:Username/photos/:PhotoId/", rt.wrap(rt.deletePhoto))

	// likePhoto //
	rt.router.POST("/users/:Username/photos/:PhotoId/likes/", rt.wrap(rt.likePhoto))

	// unlikePhoto //
	rt.router.DELETE("/users/:Username/photos/:PhotoId/likes/:LikeId", rt.wrap(rt.unlikePhoto))

	// commentPhoto //
	rt.router.POST("/users/:Username/photos/:PhotoId/comments/", rt.wrap(rt.commentPhoto))

	// uncommentPhoto //
	rt.router.DELETE("/users/:Username/photos/:PhotoId/comments/:CommentId", rt.wrap(rt.uncommentPhoto))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
