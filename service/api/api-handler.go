package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.wrap(rt.doLogin))
	rt.router.POST("/users/", rt.wrap(rt.addUser))
	rt.router.GET("/users/", rt.wrap(rt.getUserProfile))
	rt.router.DELETE("/users/:Username/", rt.wrap(rt.deleteUser))
	rt.router.PUT("/users/:Username/username/", rt.wrap(rt.setMyUsername))
	rt.router.POST("/users/:Username/followers/", rt.wrap(rt.followUser))
	rt.router.DELETE("/users/:Username/followers/:FollowerId", rt.wrap(rt.unfollowUser))
	rt.router.POST("/users/:Username/bannedusers/", rt.wrap(rt.banUser))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
