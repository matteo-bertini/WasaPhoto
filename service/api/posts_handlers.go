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

// supportedImageExts maps each accepted file extension to its HTTP Content-Type,
// kept in sync with the formats validated in UploadPostHandler (JPEG, PNG).
var supportedImageExts = []struct {
	ext         string
	contentType string
}{
	{".jpg", "image/jpeg"},
	{".png", "image/png"},
}

// locatePostImage finds the on-disk file for a post, trying every supported
// extension, since the original upload format (JPEG/PNG) is not tracked elsewhere.
func locatePostImage(authorID, postID string) (filePath string, contentType string, found bool) {
	for _, f := range supportedImageExts {
		candidate := filepath.Join(UploadDir, authorID, postID+f.ext)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, f.contentType, true
		}
	}
	return "", "", false
}

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

	// 5-bis. Detect the real content type from the file bytes (never trust the
	// client-supplied filename/extension) and only accept JPEG/PNG, as declared
	// in the spec. Rewind the reader afterwards so io.Copy below gets the full file.
	sniff := make([]byte, 512)
	n, err := file.Read(sniff)
	if err != nil && !errors.Is(err, io.EOF) {
		ctx.Logger.WithError(err).Error("UploadPostHandler: failed to read file for content-type detection")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var fileExt string
	switch http.DetectContentType(sniff[:n]) {
	case "image/jpeg":
		fileExt = ".jpg"
	case "image/png":
		fileExt = ".png"
	default:
		ctx.Logger.Warn("UploadPostHandler: unsupported file format, only JPEG/PNG are allowed")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		ctx.Logger.WithError(err).Error("UploadPostHandler: failed to rewind uploaded file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	caption := r.FormValue("caption")

	if len(caption) > 2200 {
		ctx.Logger.Error("UploadPostHandler: caption exceeds 2200 characters")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 6. Generate a KSUID for the post ID
	id, err := ksuid.NewRandom()
	if err != nil {
		ctx.Logger.WithError(err).Error("UploadPostHandler: failed to generate postID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postID := id.String()

	// 7. Define and ensure the user-specific upload directory exists
	userUploadDir := filepath.Join(UploadDir, ctx.UserID)

	// os.MkdirAll is efficient: it creates the path if it doesn't exist, otherwise does nothing
	if err := os.MkdirAll(userUploadDir, 0755); err != nil {
		ctx.Logger.WithError(err).WithField("path", userUploadDir).Error("UploadPostHandler: failed to create user directory")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 8. Create physical file on disk inside the user's folder, using the
	// extension that matches the detected content type (.jpg or .png)
	filePath := filepath.Join(userUploadDir, postID+fileExt)
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

// GetPhotoHandler retrieves the image from the user-specific directory using the username in the path.
func (rt *_router) GetPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// 1. Extract parameters from the path
	pathUsername := ps.ByName("username")
	postID := ps.ByName("postId")

	// 2. Translate username to userID to find the correct folder
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

	// 3. Locate the file: it was stored as either postID.jpg or postID.png
	// depending on the format detected at upload time.
	filePath, contentType, found := locatePostImage(authorID, postID)
	if !found {
		ctx.Logger.WithField("postID", postID).Warn("GetPhoto: image file not found on disk")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// 4. Success: Serve the binary data with the matching Content-Type
	w.Header().Set("Content-Type", contentType)
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

	// 2. Physical File Cleanup: the file may be postID.jpg or postID.png
	if filePath, _, found := locatePostImage(loggedUserID, postID); found {
		if err := os.Remove(filePath); err != nil {
			ctx.Logger.WithError(err).WithField("path", filePath).Warn("DeletePostHandler: could not delete file from disk")
		}
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
	if err := json.NewEncoder(w).Encode(likes); err != nil {
		// The 200 status and part of the body may have already been flushed by
		// Encode, so calling w.WriteHeader here would just be a superfluous/
		// ineffective call. We only log the failure.
		ctx.Logger.WithError(err).Error("GetLikesHandler: failed to encode response")
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

	// Enforce the maxLength: 500 constraint declared on CommentContent.content in doc/api.yaml
	if len(body.Content) > 500 {
		ctx.Logger.Warn("AddCommentHandler: content exceeds 500 characters")
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
