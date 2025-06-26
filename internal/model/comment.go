package model

import (
	"database/sql"
	_ "embed"
	"errors"

	"github.com/marcusprice/twitter-clone/internal/dbutils"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/logger"
)

type CommentModel struct {
	db *sql.DB
}

//go:embed queries/select-comment-by-id.sql
var selectCommentByIDQuery string

func (commentModel *CommentModel) GetByID(postID int) (dtypes.CommentData, error) {
	var id int
	var post_id int
	var user_id int
	var depth int
	var parent_comment_id sql.NullInt64
	var content string
	var image string
	var like_count int
	var retweet_count int
	var bookmark_count int
	var impressions int
	var created_at string
	var updated_at string
	var author_username string
	var author_display_name string
	var author_avatar string

	err := commentModel.db.
		QueryRow(selectCommentByIDQuery, postID).
		Scan(
			&id, &post_id, &user_id, &depth, &parent_comment_id, &content,
			&image, &like_count, &retweet_count, &bookmark_count, &impressions,
			&created_at, &updated_at, &author_username, &author_display_name,
			&author_avatar)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dtypes.CommentData{}, CommentNotFoundError{}
		}

		return dtypes.CommentData{}, err
	}

	author := dtypes.Author{
		Username:    author_username,
		DisplayName: author_display_name,
		Avatar:      author_avatar,
	}

	commentData := dtypes.CommentData{
		ID:              id,
		PostID:          post_id,
		UserID:          user_id,
		Depth:           depth,
		ParentCommentID: int(parent_comment_id.Int64),
		Content:         content,
		Image:           image,
		LikeCount:       like_count,
		RetweetCount:    retweet_count,
		BookmarkCount:   bookmark_count,
		Impressions:     impressions,
		CreatedAt:       created_at,
		UpdatedAt:       updated_at,
		Author:          author,
	}

	return commentData, nil
}

//go:embed queries/create-post-comment.sql
var creatPostCommentQuery string

func (commentModel *CommentModel) NewPostComment(commentInput dtypes.CommentInput) (rowID int, err error) {
	err = commentModel.db.
		QueryRow(
			creatPostCommentQuery, commentInput.UserID, commentInput.PostID,
			commentInput.Content, commentInput.Image).
		Scan(&rowID)

	if err != nil {
		logger.LogError("CommentModel.NewPostComment() error: " + err.Error())
		if dbutils.ConstraintFailed(err) {
			return -1, dbutils.WrapConstraintError(err)
		}

		return -1, err
	}

	return rowID, nil
}

//go:embed queries/create-comment-reply.sql
var creatCommentReplyQuery string

func (commentModel *CommentModel) NewCommentReply(commentInput dtypes.CommentInput) (rowID int, err error) {
	err = commentModel.db.
		QueryRow(
			creatCommentReplyQuery, commentInput.UserID, commentInput.PostID,
			commentInput.ParentCommentID, commentInput.Content,
			commentInput.Image).
		Scan(&rowID)

	if err != nil {
		logger.LogError("CommentModel.NewPostComment() error: " + err.Error())
		if dbutils.IsConstraintError(err) {
			return -1, dbutils.WrapConstraintError(err)
		}

		return -1, err
	}

	return rowID, nil

}

func NewCommentModel(db *sql.DB) *CommentModel {
	return &CommentModel{db}
}
