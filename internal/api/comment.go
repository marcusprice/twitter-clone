package api

import (
	"net/http"

	"github.com/marcusprice/twitter-clone/internal/controller"
)

type CommentAPI struct {
	comment *controller.Comment
}

func (commentAPI *CommentAPI) Create(w http.ResponseWriter, r *http.Request) {

}

func NewCommentAPI(comment *controller.Comment) *CommentAPI {
	return &CommentAPI{comment: comment}
}
