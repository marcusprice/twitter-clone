package controller

import (
	"database/sql"
	"errors"

	"github.com/marcusprice/twitter-clone/internal/model"
)

type Timeline struct {
	userID    int
	userModel *model.UserModel
	postModel *model.PostModel
	posts     []*Post
}

func (t *Timeline) SetID(userID int) *Timeline {
	t.userID = userID
	return t
}

func (t *Timeline) GetPosts(limit, offset int) (posts []*Post, postsRemaining int, err error) {
	if t.userID == 0 {
		return []*Post{}, -1, errors.New("userID required to fetch posts")
	}

	postRows, postIDs, err := t.postModel.QueryUserTimeline(t.userID, limit, offset)
	if err != nil {
		return []*Post{}, -1, err
	}

	postsRemaining, err = t.postModel.TimelineRemainingPostsCount(t.userID, limit, offset)
	rowsAffected, _ := t.postModel.AddImpressionBulk(postIDs) // okay to silently fail

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

	return posts, postsRemaining, nil
}

func NewTimelineController(db *sql.DB) *Timeline {
	return &Timeline{
		userModel: model.NewUserModel(db),
		postModel: model.NewPostModel(db),
	}
}
