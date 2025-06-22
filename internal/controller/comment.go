package controller

import (
	"database/sql"

	"github.com/marcusprice/twitter-clone/internal/model"
)

type Comment struct {
	commentModel *model.CommentModel
}

func NewCommentController(db *sql.DB) *Comment {
	commentModel := model.NewCommentModel(db)
	return &Comment{commentModel}
}
