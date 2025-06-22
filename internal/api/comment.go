package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/marcusprice/twitter-clone/internal/controller"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
)

type CommentAPI struct {
	comment *controller.Comment
}

func (commentAPI *CommentAPI) Create(w http.ResponseWriter, r *http.Request) {
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

	postIDFormValue := r.FormValue("postID")
	parentCommentIDFormValue := r.FormValue("parentCommentID")
	content = r.FormValue("content")

	postID, err := strconv.Atoi(postIDFormValue)
	if err != nil {
		http.Error(w, BadRequest, http.StatusBadRequest)
		return
	}

	parentCommentID := 0
	if parentCommentIDFormValue != "" {
		parentCommentID, err = strconv.Atoi(parentCommentIDFormValue)
		if err != nil {
			http.Error(w, BadRequest, http.StatusBadRequest)
			return
		}
	}

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

	commentInput := dtypes.CommentInput{
		UserID:          userID,
		PostID:          postID,
		ParentCommentID: parentCommentID,
		Content:         content,
		Image:           filename,
	}

	comment, err := commentAPI.comment.New(commentInput)
	if err != nil {
		if errors.Is(err, controller.DepthLimitError{}) {
			http.Error(w, BadRequest, http.StatusBadRequest)
		} else {
			http.Error(w, InternalServerError, http.StatusInternalServerError)
		}

		return
	}

	payload := generateCommentPayload(comment)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func NewCommentAPI(comment *controller.Comment) *CommentAPI {
	return &CommentAPI{comment: comment}
}
