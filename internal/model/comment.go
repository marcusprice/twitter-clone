package model

import "database/sql"

type CommentModel struct {
	db *sql.DB
}

func NewCommentModel(db *sql.DB) *CommentModel {
	return &CommentModel{db}
}
