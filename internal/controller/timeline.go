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

	postRows, err := t.postModel.QueryUserTimeline(t.userID, limit, offset)
	if err != nil {
		return []*Post{}, -1, err
	}

	posts = []*Post{}

	for _, row := range postRows {
		post := &Post{}
		post.setFromModel(row)
		post.AddImpression()
		posts = append(posts, post)
	}

	postsRemaining, err = t.postModel.TimelineRemainingPostsCount(
		t.userID, limit, offset,
	)

	if err != nil {
		return []*Post{}, -1, err
	}

	return posts, postsRemaining - limit, nil
}

func NewTimelineController(db *sql.DB) *Timeline {
	return &Timeline{
		userModel: model.NewUserModel(db),
		postModel: model.NewPostModel(db),
	}
}
