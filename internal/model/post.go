package model

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"

	"github.com/marcusprice/twitter-clone/internal/dbutils"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/logger"
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

func (pm PostModel) GetByID(id int) (dtypes.PostData, error) {
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
			return dtypes.PostData{}, PostNotFoundError{}
		} else {
			if util.InDevContext() {
				panic(err)
			}
			return dtypes.PostData{}, err
		}
	}

	postAuthor := dtypes.PostAuthor{
		Username:    username,
		DisplayName: displayName,
		Avatar:      avatar,
	}

	postData := dtypes.PostData{
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

//go:embed queries/user-timeline-query.sql
var userTimelineBaseQuery string

func (post *PostModel) QueryUserTimeline(userID, limit, offset int) ([]dtypes.PostData, error) {
	if limit <= 0 {
		logger.LogError("PostModel.TimelineRemainingPostsCount(): postitive limit value required")
		return []dtypes.PostData{}, errors.New("Positive limit value required")
	}

	result, err := post.db.Query(userTimelineBaseQuery, userID, limit, offset)
	if err != nil {
		logger.LogError("PostModel.QueryUserTimeline(): query error: " + err.Error())
		return []dtypes.PostData{}, err
	}

	ret := []dtypes.PostData{}
	for result.Next() {
		var id int
		var user_id int
		var content string
		var like_count int
		var retweet_count int
		var bookmark_count int
		var impressions int
		var image string
		var created_at string
		var updated_at string
		var user_name string
		var display_name string
		var avatar string

		err := result.Scan(
			&id, &user_id, &content, &like_count, &retweet_count,
			&bookmark_count, &impressions, &image, &created_at, &updated_at,
			&user_name, &display_name, &avatar)

		if err != nil {
			logger.LogError("PostModel.QueryUserTimeline(): error scanning timeline post: " + err.Error())
			return []dtypes.PostData{}, err
		}

		postAuthor := dtypes.PostAuthor{
			Username:    user_name,
			DisplayName: display_name,
			Avatar:      avatar,
		}

		postData := dtypes.PostData{
			ID:            id,
			UserID:        user_id,
			Content:       content,
			LikeCount:     like_count,
			RetweetCount:  retweet_count,
			BookmarkCount: bookmark_count,
			Impressions:   impressions,
			Image:         image,
			CreatedAt:     created_at,
			UpdatedAt:     updated_at,
			Author:        postAuthor,
		}

		ret = append(ret, postData)
	}

	return ret, nil
}

//go:embed queries/timeline-offset-count.sql
var timelineOffsetCountQuery string

func (pm *PostModel) TimelineRemainingPostsCount(userID, limit, offset int) (int, error) {
	if limit <= 0 {
		logger.LogError("PostModel.TimelineRemainingPostsCount(): postitive limit value required")
		return -1, errors.New("Positive limit value required")
	}

	var count int

	err := pm.db.
		QueryRow(timelineOffsetCountQuery, userID, limit, offset).
		Scan(&count)

	if err != nil {
		logger.LogError("PostModel.OffsetCount() error: " + err.Error())
		return -1, err
	}

	remainingPosts := count - (limit + offset)
	if remainingPosts < 0 {
		remainingPosts = 0
	}

	return remainingPosts, nil
}

//go:embed queries/add-impression.sql
var addImpressionQuery string

func (postModel *PostModel) AddImpression(postID int) error {
	result, err := postModel.db.Exec(addImpressionQuery, postID)
	if err != nil {
		logMsg := fmt.Sprintf(
			"Error adding impression for postID: %d, error: %s",
			postID, err.Error())

		logger.LogError(logMsg)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logMsg := fmt.Sprintf(
			"Error adding impression for postID: %d, error: %s",
			postID, err.Error())

		logger.LogError(logMsg)
	}

	if rowsAffected == 0 {
		logMsg := fmt.Sprintf("No rows affected for postID: %d", postID)
		logger.LogError(logMsg)
	}

	return nil
}

func NewPostModel(db *sql.DB) *PostModel {
	return &PostModel{db}
}
