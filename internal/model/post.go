package model

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strconv"

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
	var comment_count int
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
			&comment_count, &likeCount, &retweetCount, &bookmarkCount,
			&impressions, &image, &createdAt, &updatedAt)

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

	postAuthor := dtypes.Author{
		Username:    username,
		DisplayName: displayName,
		Avatar:      avatar,
	}

	postData := dtypes.PostData{
		Author:        postAuthor,
		ID:            postID,
		UserID:        userID,
		Content:       content,
		CommentCount:  comment_count,
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

//go:embed queries/select-post-by-id-user-context.sql
var selectByIDUserContextQuery string

func (pm PostModel) GetByIDUserContext(userID, postID int) (dtypes.PostData, error) {
	var username string
	var displayName string
	var avatar string
	var id int
	var authorID int
	var content string
	var comment_count int
	var likeCount int
	var retweetCount int
	var bookmarkCount int
	var impressions int
	var image string
	var createdAt string
	var updatedAt string
	var liked int

	err := pm.db.
		QueryRow(selectByIDUserContextQuery, userID, postID).
		Scan(
			&username, &displayName, &avatar, &id, &authorID, &content,
			&comment_count, &likeCount, &retweetCount, &bookmarkCount,
			&impressions, &image, &createdAt, &updatedAt, &liked)

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

	postAuthor := dtypes.Author{
		Username:    username,
		DisplayName: displayName,
		Avatar:      avatar,
	}

	postData := dtypes.PostData{
		Author:        postAuthor,
		ID:            id,
		UserID:        authorID,
		Content:       content,
		CommentCount:  comment_count,
		LikeCount:     likeCount,
		RetweetCount:  retweetCount,
		BookmarkCount: bookmarkCount,
		Impressions:   impressions,
		Image:         image,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
		Liked:         liked,
	}

	return postData, nil
}

//go:embed queries/user-timeline-query.sql
var userTimelineQuery string

func (post *PostModel) QueryUserFollowingTimeline(userID, limit, offset int) (postRows []dtypes.TimelinePostData, postIDs []int, err error) {
	if limit <= 0 {
		logger.LogError("PostModel.TimelineRemainingPostsCount(): postitive limit value required")
		return []dtypes.TimelinePostData{}, []int{}, errors.New("Positive limit value required")
	}

	result, err := post.db.Query(userTimelineQuery, userID, limit, offset)
	if err != nil {
		logger.LogError("PostModel.QueryUserTimeline(): query error: " + err.Error())
		return []dtypes.TimelinePostData{}, []int{}, err
	}

	for result.Next() {
		postData, postID, err := parseTimelineRow(result)

		if err != nil {
			return []dtypes.TimelinePostData{}, []int{}, err
		}

		if postData.Type != "comment-retweet" {
			postIDs = append(postIDs, postID)
		}
		postRows = append(postRows, postData)
	}

	return postRows, postIDs, nil
}

//go:embed queries/user-timeline-query-all-posts.sql
var selectAllPostsQuery string

func (pm *PostModel) GetAllIncludingRetweets(userID, limit, offset int) (postRows []dtypes.TimelinePostData, postIDs []int, err error) {
	result, err := pm.db.Query(selectAllPostsQuery, userID, limit, offset)
	if err != nil {
		// TODO handle error
		logger.LogError("PostModel.GetAll() - error querying posts: " + err.Error())
	}

	for result.Next() {
		postData, postID, err := parseTimelineRow(result)

		if err != nil {
			return []dtypes.TimelinePostData{}, []int{}, err
		}

		if postData.Type != "comment-retweet" {
			postIDs = append(postIDs, postID)
		}

		postRows = append(postRows, postData)
	}

	return postRows, postIDs, nil
}

//go:embed queries/user-timeline-all-posts-count.sql
var userTimelineAllPostsCount string

func (pm *PostModel) AllIncludingRetweetCount() (int, error) {
	var count int
	err := pm.db.QueryRow(userTimelineAllPostsCount).Scan(&count)
	if err != nil {
		logger.LogError("PostModel.GetAllCount() - error scanning row: " + err.Error())
		return -1, err
	}

	return count, nil
}

//go:embed queries/user-timeline-following-count.sql
var timelineOffsetCountQuery string

func (pm *PostModel) UserFollowingTimelineCount(userID int) (int, error) {
	var count int

	err := pm.db.QueryRow(timelineOffsetCountQuery, userID).Scan(&count)
	if err != nil {
		logger.LogError("PostModel.OffsetCount() - error: " + err.Error())
		return -1, err
	}

	return count, nil
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

//go:embed queries/add-impression-bulk.sql
var addImpressionBulkQuery string

func (postModel *PostModel) AddImpressionBulk(postIDs []int) (rowsAffected int, err error) {
	var idStrings string
	for index, id := range postIDs {
		if index == 0 {
			idStrings += strconv.Itoa(id)
		} else {
			idStrings += fmt.Sprintf(", %d", id)
		}
	}

	query := fmt.Sprintf(addImpressionBulkQuery, idStrings)

	result, err := postModel.db.Exec(query)
	if err != nil {
		return -1, err
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return -1, err
	}

	if ra == 0 {
		logger.LogError("PostModel.AddImpressionBulk(): no rows affected")
		return -1, errors.New("no rows affected")
	}

	return int(ra), nil
}

func parseTimelineRow(result dbutils.RowScanner) (postData dtypes.TimelinePostData, postID int, err error) {
	var content_type string
	var id int
	var user_id int
	var content string
	var comment_count int
	var like_count int
	var retweet_count int
	var bookmark_count int
	var impressions int
	var image string
	var created_at string
	var updated_at string
	var author_user_name string
	var author_display_name string
	var author_avatar string
	var parent_post_id sql.NullInt64
	var parent_post_author_username sql.NullString
	var parent_comment_id sql.NullInt64
	var parent_comment_author_username sql.NullString
	var retweeter_user_name_ns sql.NullString
	var retweeter_display_name_ns sql.NullString
	var liked int
	var retweeted int
	var bookmarked int
	var sort_throwaway string

	err = result.Scan(
		&content_type, &id, &user_id, &content, &comment_count, &like_count,
		&retweet_count, &bookmark_count, &impressions, &image, &created_at,
		&updated_at, &author_user_name, &author_display_name, &author_avatar,
		&retweeter_user_name_ns, &retweeter_display_name_ns, &parent_post_id,
		&parent_post_author_username, &parent_comment_id,
		&parent_comment_author_username, &liked, &retweeted, &bookmarked,
		&sort_throwaway)

	if err != nil {
		logger.LogError("PostModel.parseTimelineRow(): error scanning timeline post: " + err.Error())
		return dtypes.TimelinePostData{}, -1, err
	}

	postRetweeter := dtypes.Retweeter{
		Username:    retweeter_user_name_ns.String,
		DisplayName: retweeter_display_name_ns.String,
	}

	postAuthor := dtypes.Author{
		Username:    author_user_name,
		DisplayName: author_display_name,
		Avatar:      author_avatar,
	}

	postData = dtypes.TimelinePostData{
		Type:                        content_type,
		ID:                          id,
		UserID:                      user_id,
		Content:                     content,
		CommentCount:                comment_count,
		LikeCount:                   like_count,
		RetweetCount:                retweet_count,
		BookmarkCount:               bookmark_count,
		Impressions:                 impressions,
		Image:                       image,
		CreatedAt:                   created_at,
		UpdatedAt:                   updated_at,
		ViewerLiked:                 liked,
		ViewerRetweeted:             retweeted,
		ViewerBookmarked:            bookmarked,
		ParentPostID:                int(parent_post_id.Int64),
		ParentPostAuthorUsername:    parent_post_author_username.String,
		ParentCommentID:             int(parent_comment_id.Int64),
		ParentCommentAuthorUsername: parent_comment_author_username.String,

		Author:    postAuthor,
		Retweeter: postRetweeter,
	}

	return postData, id, nil
}

func NewPostModel(db *sql.DB) *PostModel {
	return &PostModel{db}
}
