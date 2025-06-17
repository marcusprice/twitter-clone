package model

import (
	"database/sql"
	_ "embed"

	"github.com/marcusprice/twitter-clone/internal/dbutils"
)

type PostAction struct {
	db *sql.DB
}

//go:embed queries/delete-post-like.sql
var deletePostLikeQuery string

//go:embed queries/create-post-like.sql
var createPostLikeQuery string

func (pa *PostAction) Like(postID, userID int) error {
	result, err := pa.db.Exec(createPostLikeQuery, postID, userID)
	if err != nil {
		if dbutils.IsUniqueConstraintError(err) {
			// user already likes this post, likely a duplicate request
			return nil
		}

		if dbutils.ConstraintFailed(err) {
			return dbutils.WrapConstraintError(err)
		}

		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || int(rowsAffected) == 0 {
		return err
	}

	return nil
}

func (pa *PostAction) Unlike(postID, userID int) error {
	result, err := pa.db.Exec(deletePostLikeQuery, postID, userID)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func NewPostActionModel(db *sql.DB) *PostAction {
	return &PostAction{db}
}
