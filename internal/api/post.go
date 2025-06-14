package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/dbutils"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
)

const MAX_POST_UPLOAD_BYTES int64 = 1024 * 1024 * 10 // 10 mb

type PostAPI struct {
	post *controller.Post
}

func (postAPI PostAPI) CreatePost(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		statusText := http.StatusText(http.StatusInternalServerError)
		http.Error(w, statusText, http.StatusInternalServerError)
		return
	}

	filename := ""
	content := ""

	r.Body = http.MaxBytesReader(w, r.Body, MAX_POST_UPLOAD_BYTES)
	err := r.ParseMultipartForm(getMaxUploadMemory())
	if err != nil {
		if requestBodyTooLarge(err) {
			statusText := http.StatusText(http.StatusRequestEntityTooLarge)
			http.Error(w, statusText, http.StatusRequestEntityTooLarge)
		} else {
			statusText := http.StatusText(http.StatusInternalServerError)
			http.Error(w, statusText, http.StatusInternalServerError)
		}

		return
	}

	content = r.FormValue("content")
	file, header, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		statusText := http.StatusText(http.StatusBadRequest)
		http.Error(w, statusText, http.StatusBadRequest)
		return
	}

	if err == http.ErrMissingFile {
		// no file upload, content required
		if content == "" {
			statusText := http.StatusText(http.StatusBadRequest)
			http.Error(w, statusText, http.StatusBadRequest)
			return
		}
	} else {
		// user uploaded file, content optional
		defer file.Close()

		filename, err = handleImageUpload(file, header)
		if err != nil {
			var invalidFileTypeError InvalidFileTypeError
			if errors.As(err, &invalidFileTypeError) {
				statusText := http.StatusText(http.StatusUnsupportedMediaType)
				http.Error(w, statusText, http.StatusUnsupportedMediaType)
			} else {
				statusText := http.StatusText(http.StatusInternalServerError)
				http.Error(w, statusText, http.StatusInternalServerError)
			}

			return
		}
	}

	postInput := dtypes.PostInput{
		UserID:  userID,
		Content: content,
		Image:   filename,
	}

	err = postAPI.post.New(postInput)
	if err != nil {
		if dbutils.IsConstraintError(err) {
			statusText := http.StatusText(http.StatusBadRequest)
			http.Error(w, statusText, http.StatusBadRequest)
		} else {
			statusText := http.StatusText(http.StatusInternalServerError)
			http.Error(w, statusText, http.StatusInternalServerError)
		}

		return
	}

	payload := generatePostPayload(postAPI.post)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func generatePostPayload(post *controller.Post) PostPayload {
	author := AuthorPayload{
		Username:    post.Author.Username,
		DisplayName: post.Author.DisplayName,
		Avatar:      post.Author.Avatar,
	}

	return PostPayload{
		ID:            post.ID,
		Content:       post.Content,
		LikeCount:     post.LikeCount,
		RetweetCount:  post.RetweetCount,
		BookmarkCount: post.BookmarkCount,
		Impressions:   post.Impressions,
		Image:         post.Image,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
		Author:        author,
	}
}

func NewPostAPI(postController *controller.Post) *PostAPI {
	return &PostAPI{postController}
}
