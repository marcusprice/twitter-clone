package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/dbutils"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/model"
)

const MAX_POST_UPLOAD_BYTES int64 = 1024 * 1024 * 10 // 10 mb

type PostAPI struct {
	post *controller.Post
}

func (postAPI PostAPI) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	postIDPathValue := r.PathValue("postID")
	postID, err := strconv.Atoi(postIDPathValue)
	if err != nil {
		http.Error(w, BadRequest, http.StatusBadRequest)
		return
	}

	post, err := postAPI.post.GetPostAndComments(postID, userID)
	if err != nil {
		// TODO: figure out error handling
		http.Error(w, InternalServerError, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(generatePostAndCommentsPayload(post))
}

func (postAPI PostAPI) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	filename := ""
	content := ""

	r.Body = http.MaxBytesReader(w, r.Body, MAX_POST_UPLOAD_BYTES)
	err := r.ParseMultipartForm(getMaxUploadMemory())
	if err != nil {
		if requestBodyTooLarge(err) {
			http.Error(w, RequestEntityTooLarge, http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, InternalServerError, http.StatusInternalServerError)
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
			http.Error(w, BadRequest, http.StatusBadRequest)
			return
		}
	} else {
		// user uploaded file, content optional
		defer file.Close()

		filename, err = handleImageUpload(file, header)
		if err != nil {
			var invalidFileTypeError InvalidFileTypeError
			if errors.As(err, &invalidFileTypeError) {
				http.Error(w, UnsupportedMediaType, http.StatusUnsupportedMediaType)
			} else {
				http.Error(w, InternalServerError, http.StatusInternalServerError)
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
			http.Error(w, BadRequest, http.StatusBadRequest)
		} else {
			http.Error(w, InternalServerError, http.StatusInternalServerError)
		}

		return
	}

	payload := generatePostPayload(postAPI.post)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func (postAPI *PostAPI) Like(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, BadRequest, http.StatusBadRequest)
	}

	post := postAPI.post
	err = post.ByID(postID)
	if err != nil {
		var postNotFoundError model.PostNotFoundError
		if errors.As(err, &postNotFoundError) {
			http.Error(w, NotFound, http.StatusNotFound)
		} else {
			http.Error(w, InternalServerError, http.StatusInternalServerError)
		}
	}

	if r.Method == http.MethodPut {
		err = post.Like(userID)
	} else {
		err = post.Unlike(userID)
	}

	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (postAPI *PostAPI) Retweet(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, BadRequest, http.StatusBadRequest)
	}

	post := postAPI.post
	err = post.ByID(postID)
	if err != nil {
		var postNotFoundError model.PostNotFoundError
		if errors.As(err, &postNotFoundError) {
			http.Error(w, NotFound, http.StatusNotFound)
		} else {
			http.Error(w, InternalServerError, http.StatusInternalServerError)
		}
	}

	if r.Method == http.MethodPut {
		err = post.Retweet(userID)
	} else {
		err = post.UnRetweet(userID)
	}

	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (postAPI *PostAPI) Bookmark(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	postID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, BadRequest, http.StatusBadRequest)
	}

	post := postAPI.post
	err = post.ByID(postID)
	if err != nil {
		var postNotFoundError model.PostNotFoundError
		if errors.As(err, &postNotFoundError) {
			http.Error(w, NotFound, http.StatusNotFound)
		} else {
			http.Error(w, InternalServerError, http.StatusInternalServerError)
		}
	}

	if r.Method == http.MethodPut {
		err = post.Bookmark(userID)
	} else {
		err = post.UnBookmark(userID)
	}

	if err != nil {
		http.Error(w, InternalServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func NewPostAPI(postController *controller.Post) *PostAPI {
	return &PostAPI{postController}
}
