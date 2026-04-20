package api

import "WasaPhoto/service/database"

// setMyUsername operation //
type setMyUsernameRequestBody struct {
	Username string `json:"Username"`
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
