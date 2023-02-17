package api

import "WasaPhoto/service/database"

// In questo file vengono dichiarate le struct necessarie per il funzionamento del programma //

type Photo struct {
	PhotoId        string `json:"PhotoId"`
	LikesNumber    int    `json:"LikesNumber"`
	CommentsNumber int    `json:"CommentsNumber"`
	DateOfUpload   string `json:"DateOfUpload"`
}

//---------------------------------- RequestBody e ResponseBody delle varie operazioni -----------------------------------------//

// doLogin operation //
type doLoginRequestBody struct {
	Username string `json:"Username"`
}
type doLoginResponseBody struct {
	Identifier string `json:"Identifier"`
}

// addUser operation //
type addUserRequestBody struct {
	Username string `json:"Username"`
}
type addUserResponseBody struct {
	Username       string  `json:"Username"`
	Followers      int     `json:"Followers"`
	Following      int     `json:"Following"`
	NumberOfPhotos int     `json:"NumberOfPhotos"`
	UploadedPhotos []Photo `json:"UploadedPhotos"`
}

// getUserProfile operation //
type getUserProfileResponseBody struct {
	Username       string  `json:"Username"`
	Followers      int     `json:"Followers"`
	Following      int     `json:"Following"`
	NumberOfPhotos int     `json:"NumberOfPhotos"`
	UploadedPhotos []Photo `json:"UploadedPhotos"`
}

func (r *getUserProfileResponseBody) FromDatabase(db_user database.Database_user) {
	r.Username = db_user.Username
	r.Followers = db_user.Followers
	r.Following = db_user.Following
	r.NumberOfPhotos = db_user.Numberofphotos
	UploadedPhotos := []Photo{}
	for _, ph := range db_user.UploadedPhotos {
		var to_add Photo
		to_add.PhotoId = ph.PhotoId
		to_add.LikesNumber = ph.LikesNumber
		to_add.CommentsNumber = ph.CommentsNumber
		to_add.DateOfUpload = ph.DateOfUpload
		UploadedPhotos = append(UploadedPhotos, to_add)
	}
	r.UploadedPhotos = UploadedPhotos
	return

}

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

// uploadPhoto operation //
func (ph *Photo) PhotoToDatabase() database.Database_photo {
	var to_ret database.Database_photo
	to_ret.PhotoId = ph.PhotoId
	to_ret.LikesNumber = ph.LikesNumber
	to_ret.CommentsNumber = ph.CommentsNumber
	to_ret.DateOfUpload = ph.DateOfUpload
	return to_ret
}

// likePhoto Operation //
// In questo caso il RequestBody è uguale al ResponseBody in caso di successo //
type likePhotoResponseBody struct {
	LikeId string `json:"LikeId"`
}

type commentPhotoRequestBody struct {
	CommentId     string `json:"CommentId"`
	CommentAuthor string `json:"CommentAuthor"`
	CommentText   string `json:"CommentText"`
}
