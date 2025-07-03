package controller

import (
	"database/sql"
	"errors"

	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/model"
)

type TimelineView string

const FOLLOWING TimelineView = "FOLLOWING"
const FOR_YOU TimelineView = "FOR_YOU"

var TIMELINE_VIEWS []TimelineView = []TimelineView{FOLLOWING, FOR_YOU}

type Timeline struct {
	userID    int
	view      TimelineView
	userModel *model.UserModel
	postModel *model.PostModel
	posts     []*Post
}

func (t *Timeline) Set(userID int, view TimelineView) *Timeline {
	t.userID = userID
	t.view = view
	return t
}

func (t *Timeline) GetPosts(limit, offset int) (posts []*Post, postsRemaining int, err error) {
	if t.userID == 0 {
		return []*Post{}, -1, errors.New("userID required to fetch posts")
	}

	var postRows []dtypes.PostData
	var postIDs []int
	var totalPosts int
	if t.view == FOLLOWING {
		postRows, postIDs, err = t.postModel.QueryUserFollowingTimeline(t.userID, limit, offset)
		if err != nil {
			return []*Post{}, -1, err
		}

		totalPosts, err = t.postModel.UserFollowingTimelineCount(t.userID)
		if err != nil {
			return []*Post{}, -1, err
		}
	} else {
		postRows, postIDs, err = t.postModel.GetAllIncludingRetweets(t.userID, limit, offset)
		if err != nil {
			return []*Post{}, -1, err
		}

		totalPosts, err = t.postModel.AllIncludingRetweetCount()
		if err != nil {
			return []*Post{}, -1, err
		}
	}

	rowsAffected := 0
	// TODO: this is a performance bottleneck
	if len(postRows) > 0 {
		rowsAffected, _ = t.postModel.AddImpressionBulk(postIDs) // okay to silently fail
	}

	if err != nil {
		return []*Post{}, -1, err
	}

	posts = []*Post{}
	for _, row := range postRows {
		post := &Post{}
		post.setFromModel(row)
		if rowsAffected == len(postRows) {
			post.Impressions += 1
		}
		posts = append(posts, post)
		postIDs = append(postIDs, post.ID)
	}

	return posts, totalPosts - (limit + offset), nil
}

func NewTimelineController(db *sql.DB) *Timeline {
	return &Timeline{
		userModel: model.NewUserModel(db),
		postModel: model.NewPostModel(db),
	}
}
