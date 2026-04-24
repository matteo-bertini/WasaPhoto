package api

import (
	"WasaPhoto/service/api/reqcontext"
	"WasaPhoto/service/models"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
)

// UploadDir defines the physical path where images will be stored
const UploadDir = "./uploads/posts"

func (rt *_router) UploadPostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Extract the target username from the path and the authenticated UserID from context
	pathUsername := ps.ByName("username")
	loggedUserID := ctx.UserID

	// 2. Get the UserID associated with the username in the path using your existing function
	targetUserID, err := rt.db.GetIDByUsername(pathUsername)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {

			// If the user in the path doesn't exist, return 404
			ctx.Logger.WithError(models.ErrUserNotFound).Error("UploadPostHandler: target user not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// For any other DB error, return 500
		ctx.Logger.WithError(err).Error("UploadPostHandler: database error while fetching target user ID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. Authorization Check
	if targetUserID != loggedUserID {
		ctx.Logger.WithFields(logrus.Fields{
			"loggedUserID": loggedUserID,
			"targetUserID": targetUserID,
		}).Warn("UploadPostHandler: unauthorized upload attempt")
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 4. Parse MultipartForm to handle file upload (max size 10MB)
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		ctx.Logger.WithError(err).Error("UploadPostHandler: failed to parse multipart form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 5. Retrieve the file and optional caption from form data
	file, _, err := r.FormFile("file")
	if err != nil {
		ctx.Logger.WithError(err).Error("UploadPostHandler: failed to retrieve file from form")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer file.Close()

	caption := r.FormValue("caption")

	// 6. Generate a KSUID for the post ID
	id, err := ksuid.NewRandom()
	if err != nil {
		ctx.Logger.WithError(err).Error("UploadPostHandler: failed to generate postID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postID := id.String()

	// 7. Define and ensure the user-specific upload directory exists
	// Instead of a flat directory, we use UploadDir/userID/
	userUploadDir := filepath.Join(UploadDir, ctx.UserID)

	// os.MkdirAll is efficient: it creates the path if it doesn't exist, otherwise does nothing
	if err := os.MkdirAll(userUploadDir, 0755); err != nil {
		ctx.Logger.WithError(err).WithField("path", userUploadDir).Error("UploadPostHandler: failed to create user directory")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 8. Create physical file on disk inside the user's folder
	filePath := filepath.Join(userUploadDir, postID+".jpg")
	dst, err := os.Create(filePath)
	if err != nil {
		ctx.Logger.WithError(err).WithField("filePath", filePath).Error("UploadPostHandler: failed to create file on disk")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 9. Copy binary data to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		ctx.Logger.WithError(err).WithField("postID", postID).Error("UploadPostHandler: failed to copy file content")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 10. Prepare the Post object
	// Note: We don't need to store the path in DB because it's always UploadDir/authorID/postID.jpg
	newPost := models.Post{
		PostId:         postID,
		Username:       pathUsername,
		Caption:        caption,
		LikesNumber:    0,
		CommentsNumber: 0,
		DateOfUpload:   time.Now(),
		IsLikedByMe:    false,
	}

	// 11. Execute DB insertion
	err = rt.db.UploadPost(newPost)
	if err != nil {
		ctx.Logger.WithError(err).WithField("postID", postID).Error("UploadPostHandler: failed to insert post in DB")
		// Rollback: delete the file from the user's directory if DB fails
		_ = os.Remove(filePath)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 12. Final Success Response
	ctx.Logger.WithFields(logrus.Fields{
		"postID":   postID,
		"username": pathUsername,
	}).Info("UploadPostHandler: post successfully created")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newPost)
}

// GetPhoto retrieves the photo file if the requester is not banned.
// GetPhotoHandler retrieves the image from the user-specific directory using the username in the path.
func (rt *_router) GetPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// 1. Extract parameters from the path
	pathUsername := ps.ByName("username")
	postID := ps.ByName("postId")

	// 2. Translate username to userID to find the correct folder
	// We use your existing GetIDByUsername function
	authorID, err := rt.db.GetIDByUsername(pathUsername)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ctx.Logger.WithError(err).Error("GetPhoto: database error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 3. Build the correct path: UploadDir/authorID/postID.jpg
	// This matches the structure we used in UploadPost and DeleteUser
	filePath := filepath.Join(UploadDir, authorID, postID+".jpg")

	// 4. Physical existence check
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ctx.Logger.WithField("filePath", filePath).Warn("GetPhoto: image file not found on disk")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// 5. Success: Serve the binary data
	w.Header().Set("Content-Type", "image/jpeg")
	http.ServeFile(w, r, filePath)
}

func (rt *_router) DeletePostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	postID := ps.ByName("postId")
	loggedUserID := ctx.UserID

	// 1. Database deletion (checks ownership internally)
	err := rt.db.DeletePost(postID, loggedUserID)
	if err != nil {
		if errors.Is(err, models.ErrResourceNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, models.ErrForbiddenAction) {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		ctx.Logger.WithError(err).Error("DeletePostHandler: db error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 2. Physical File Cleanup
	// Path: UploadDir / loggedUserID / postID.jpg
	filePath := filepath.Join(UploadDir, loggedUserID, postID+".jpg")

	if err := os.Remove(filePath); err != nil {
		// We don't return 500 here because the DB record is already gone.
		// We just log it as a warning for future manual cleanup.
		ctx.Logger.WithError(err).WithField("path", filePath).Warn("DeletePostHandler: could not delete file from disk")
	}

	// 3. Success
	w.WriteHeader(http.StatusNoContent)
}

// LikePostHandler: PUT /users/{username}/posts/{postId}/likes
func (rt *_router) LikePostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	postID := ps.ByName("postId")
	targetUsername := ps.ByName("username")
	authUserID := ctx.UserID

	err := rt.db.LikePost(postID, authUserID, targetUsername)
	if err != nil {
		if errors.Is(err, models.ErrResourceNotFound) {
			ctx.Logger.WithField("postID", postID).Warn("LikePostHandler: post or owner not found")
			w.WriteHeader(http.StatusNotFound)
		} else if errors.Is(err, models.ErrForbiddenAction) {
			ctx.Logger.Warn("LikePostHandler: forbidden by ban")
			w.WriteHeader(http.StatusForbidden)
		} else {
			ctx.Logger.WithError(err).Error("LikePostHandler: database error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UnlikePostHandler: DELETE /users/{username}/posts/{postId}/likes
func (rt *_router) UnlikePostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	postID := ps.ByName("postId")
	authUserID := ctx.UserID

	err := rt.db.UnlikePost(postID, authUserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("UnlikePostHandler: database error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetLikesHandler retrieves the list of usernames who liked a specific post.
func (rt *_router) GetLikesHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	postID := ps.ByName("postId")

	likes, err := rt.db.GetLikes(postID, ctx.UserID)
	if err != nil {
		if errors.Is(err, models.ErrResourceNotFound) {
			ctx.Logger.WithField("postID", postID).Warn("GetLikesHandler: resource not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, models.ErrForbiddenAction) {
			ctx.Logger.WithFields(logrus.Fields{
				"postID": postID,
				"userID": ctx.UserID,
			}).Warn("GetLikesHandler: access forbidden by ban")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx.Logger.WithError(err).Error("GetLikesHandler: unexpected database error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// if likes is nil (no likes), return an empty slice [] instead of null
	if likes == nil {
		likes = []string{}
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(likes)
	if err != nil {
		ctx.Logger.WithError(err).Error("GetLikesHandler: failed to encode response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// AddCommentHandler handles POST /users/:username/posts/:postId/comments
func (rt *_router) AddCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	postID := ps.ByName("postId")

	var body struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Basic validation for empty comments
	if strings.TrimSpace(body.Content) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comment, err := rt.db.AddComment(postID, ctx.UserID, body.Content)
	if err != nil {
		if errors.Is(err, models.ErrResourceNotFound) {
			ctx.Logger.WithField("postID", postID).Warn("AddCommentHandler: post not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, models.ErrForbiddenAction) {
			ctx.Logger.Warn("AddCommentHandler: forbidden by ban")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		ctx.Logger.WithError(err).Error("AddCommentHandler: failed to add comment")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Success: return the full Comment object (201 Created)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(comment)
}
func (rt *_router) GetCommentsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	postID := ps.ByName("postId")

	comments, err := rt.db.GetComments(postID, ctx.UserID)
	if err != nil {
		if errors.Is(err, models.ErrResourceNotFound) {
			ctx.Logger.WithField("postID", postID).Warn("GetCommentsHandler: resource not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, models.ErrForbiddenAction) {
			ctx.Logger.Warn("GetCommentsHandler: forbidden by ban")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		ctx.Logger.WithError(err).Error("GetCommentsHandler: db error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(comments)
}

// DeleteCommentHandler handles DELETE /users/:username/posts/:postId/comments/:commentId
func (rt *_router) DeleteCommentHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	commentID, err := strconv.Atoi(ps.ByName("commentId"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rt.db.DeleteComment(commentID, ctx.UserID)
	if err != nil {
		if errors.Is(err, models.ErrResourceNotFound) {
			ctx.Logger.WithField("commentId", commentID).Warn("DeleteCommentHandler: comment not found")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, models.ErrForbiddenAction) {
			ctx.Logger.WithFields(logrus.Fields{
				"commentId": commentID,
				"userId":    ctx.UserID,
			}).Warn("DeleteCommentHandler: unauthorized deletion attempt")
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx.Logger.WithError(err).Error("DeleteCommentHandler: db error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
