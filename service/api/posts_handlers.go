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
		ctx.Logger.WithError(err).Error("UploadPostHandler: failed to generate KSUID")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	postID := id.String()

	// 7. Ensure the upload directory exists
	if _, err := os.Stat(UploadDir); os.IsNotExist(err) {
		err = os.MkdirAll(UploadDir, os.ModePerm)
		if err != nil {
			ctx.Logger.WithError(err).WithField("path", UploadDir).Error("UploadPostHandler: failed to create upload directory")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// 8. Create physical file on disk
	filePath := filepath.Join(UploadDir, postID+".jpg")
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

	// 10. Prepare the Post object using your specific struct
	// Note: We don't store ImagePath in the DB as per your design.
	// The image is implicitly at UploadDir/postID.jpg
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
	// Your SQL query should NOT include image_path if you removed it from the table
	err = rt.db.UploadPost(newPost)
	if err != nil {
		ctx.Logger.WithError(err).WithField("postID", postID).Error("UploadPostHandler: failed to insert post in DB")
		// Rollback: delete the file since the DB record failed
		_ = os.Remove(filePath)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 12. Final Success Response
	ctx.Logger.WithFields(logrus.Fields{
		"postID":   postID,
		"username": pathUsername,
	}).Info("UploadPostHandler: post successfully created without explicit path storage")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(newPost)
}
