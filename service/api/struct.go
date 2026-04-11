package api

import "WasaPhoto/service/database"

// setMyUsername operation //
type setMyUsernameRequestBody struct {
	Username string `json:"Username"`
}

// getFollowers operation //
type getFollowersResponseBody struct {
	Followers []database.Database_follower `json:"Followers"`
}

// getFollowing operation //
type getFollowingResponseBody struct {
	Following []database.Database_following `json:"Following"`
}

// getBanned operation //
type getBannedResponseBody struct {
	BannedUsers []database.Database_banned `json:"BannedUsers"`
}

// followUser operation //
// In questo caso il RequestBody è uguale al ResponseBody in caso di successo //
type followUserRequestBody struct {
	FollowerId string `json:"FollowerId"`
}

// banUser operation //
// In questo caso il RequestBody è uguale al ResponseBody in caso di successo //

type banUserRequestBody struct {
	BannedId string `json:"BannedId"`
}

// getComments Operation //
type getLikesResponseBody struct {
	Likes []database.Database_like `json:"Likes"`
}

// likePhoto Operation //
// In questo caso il RequestBody è uguale al ResponseBody in caso di successo //
type likePhotoResponseBody struct {
	LikeId string `json:"LikeId"`
}

// getComments Operation //
type getCommentsResponseBody struct {
	Comments []database.Database_comment `json:"Comments"`
}

// commentPhoto Operation //
type commentPhotoRequestBody struct {
	CommentId     string `json:"CommentId"`
	CommentAuthor string `json:"CommentAuthor"`
	CommentText   string `json:"CommentText"`
}

// getMyStream Operation //
type getMyStreamResponseBody struct {
	PhotoStream []database.Database_photostream_component `json:"PhotoStream"`
}
