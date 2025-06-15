package model

import (
	"database/sql"
	_ "embed"
	"errors"

	"github.com/marcusprice/twitter-clone/internal/dbutils"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/util"
)

type PostModel struct {
	db *sql.DB
}

//go:embed queries/create-post.sql
var createPostQuery string

func (pm PostModel) New(postInput dtypes.PostInput) (int, error) {
	var postID int
	err := pm.db.QueryRow(
		createPostQuery, postInput.UserID,
		postInput.Content, postInput.Image).Scan(&postID)

	if err != nil {
		if dbutils.ConstraintFailed(err) {
			return -1, dbutils.WrapConstraintError(err)
		}

		if util.InDevContext() {
			panic(err)
		}

		return -1, err
	}

	return postID, nil
}

//go:embed queries/select-post-by-id.sql
var selectPostByIdQuery string

func (pm PostModel) GetByID(id int) (PostData, error) {
	var username string
	var displayName string
	var avatar string
	var postID int
	var userID int
	var content string
	var likeCount int
	var retweetCount int
	var bookmarkCount int
	var impressions int
	var image string
	var createdAt string
	var updatedAt string

	err := pm.db.
		QueryRow(selectPostByIdQuery, id).
		Scan(
			&username, &displayName, &avatar, &postID, &userID, &content,
			&likeCount, &retweetCount, &bookmarkCount, &impressions, &image,
			&createdAt, &updatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PostData{}, PostNotFoundError{}
		} else {
			if util.InDevContext() {
				panic(err)
			}
			return PostData{}, err
		}
	}

	postAuthor := PostAuthor{
		Username:    username,
		DisplayName: displayName,
		Avatar:      avatar,
	}

	postData := PostData{
		Author:        postAuthor,
		ID:            postID,
		UserID:        userID,
		Content:       content,
		LikeCount:     likeCount,
		RetweetCount:  retweetCount,
		BookmarkCount: bookmarkCount,
		Impressions:   impressions,
		Image:         image,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	return postData, nil
}

func NewPostModel(db *sql.DB) *PostModel {
	return &PostModel{db}
}
