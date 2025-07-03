package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/logger"
	"github.com/marcusprice/twitter-clone/internal/model"
	"github.com/marcusprice/twitter-clone/internal/util"
)

type Post struct {
	model         *model.PostModel
	postAction    *model.PostAction
	comment       *Comment
	ID            int
	UserID        int
	Content       string
	CommentCount  int
	LikeCount     int
	RetweetCount  int
	BookmarkCount int
	Impressions   int
	Image         string
	Liked         bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Author        struct {
		Username    string
		DisplayName string
		Avatar      string
	}
	Comments []*Comment
	dtypes.Retweeter
}

func (p *Post) setFromModel(postData dtypes.PostData) {
	p.ID = postData.ID
	p.UserID = postData.UserID
	p.Content = postData.Content
	p.CommentCount = postData.CommentCount
	p.LikeCount = postData.LikeCount
	p.RetweetCount = postData.RetweetCount
	p.BookmarkCount = postData.BookmarkCount
	p.Impressions = postData.Impressions
	p.Image = postData.Image
	p.CreatedAt = util.ParseTime(postData.CreatedAt)
	p.UpdatedAt = util.ParseTime(postData.UpdatedAt)
	p.Author.Username = postData.Author.Username
	p.Author.DisplayName = postData.Author.DisplayName
	p.Author.Avatar = postData.Author.Avatar
	p.Retweeter.Username = postData.Retweeter.Username
	p.Retweeter.DisplayName = postData.Retweeter.DisplayName
	p.Liked = postData.Liked == 1
}

func (post *Post) New(postInput dtypes.PostInput) error {
	postID, err := post.model.New(postInput)
	if err != nil {
		return err
	}

	postData, err := post.model.GetByID(postID)
	if err != nil {
		return err
	}

	post.setFromModel(postData)

	return nil
}

func (post *Post) GetPostAndComments(postID int) (*Post, error) {
	postData, err := post.model.GetByID(postID)
	if err != nil {
		logger.LogError("Post.GetPostAndComments() error querying posts:" + err.Error())
		return &Post{}, nil
	}

	postComments, err := post.comment.GetPostComments(postData.ID)
	if err != nil {
		logger.LogError("Post.GetPostAndComments() error querying comments:" + err.Error())
		return &Post{}, nil
	}

	ret := &Post{}
	ret.setFromModel(postData)
	ret.Comments = postComments

	return ret, nil
}

// TODO: update to new patten
func (post *Post) ByID(postID int) error {
	postData, err := post.model.GetByID(postID)
	if err != nil {
		return err
	}

	post.setFromModel(postData)

	return nil
}

func (post *Post) Like(likerUserID int) error {
	if post.ID == 0 {
		err := fmt.Errorf("Post.Like(): missing required postID in post controller")
		if util.InDevContext() {
			log.Panicf("Like failed: %v", err)
		}

		return err
	}

	err := post.postAction.Like(post.ID, likerUserID)
	if err != nil {
		return err
	}

	err = post.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (post *Post) Unlike(likerUserID int) error {
	if post.ID == 0 {
		err := fmt.Errorf("Post.Unlike(): missing required postID in post controller")
		if util.InDevContext() {
			log.Panicf("Unlike failed: %v", err)
		}

		return err
	}

	err := post.postAction.Unlike(post.ID, likerUserID)
	if err != nil {
		return err
	}

	err = post.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (post *Post) Retweet(retweeterID int) error {
	if post.ID == 0 {
		err := fmt.Errorf("Post.Retweet(): missing required postID in post controller")
		if util.InDevContext() {
			log.Panicf("Like failed: %v", err)
		}

		return err
	}

	err := post.postAction.Retweet(post.ID, retweeterID)
	if err != nil {
		return err
	}

	err = post.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (post *Post) UnRetweet(retweeterID int) error {
	if post.ID == 0 {
		err := fmt.Errorf("Post.UnRetweet(): missing required postID in post controller")
		if util.InDevContext() {
			log.Panicf("Unlike failed: %v", err)
		}

		return err
	}

	err := post.postAction.UnRetweet(post.ID, retweeterID)
	if err != nil {
		return err
	}

	err = post.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (post *Post) Bookmark(bookmarkerID int) error {
	if post.ID == 0 {
		err := fmt.Errorf("Post.Bookmark(): missing required postID in post controller")
		if util.InDevContext() {
			log.Panicf("Like failed: %v", err)
		}

		return err
	}

	err := post.postAction.Bookmark(post.ID, bookmarkerID)
	if err != nil {
		return err
	}

	err = post.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (post *Post) UnBookmark(bookmarkerID int) error {
	if post.ID == 0 {
		err := fmt.Errorf("Post.UnBookmark(): missing required postID in post controller")
		if util.InDevContext() {
			log.Panicf("Unlike failed: %v", err)
		}

		return err
	}

	err := post.postAction.UnBookmark(post.ID, bookmarkerID)
	if err != nil {
		return err
	}

	err = post.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (post *Post) AddImpression() error {
	if post.ID == 0 {
		logger.LogError("Post.AddImpression(): missing postID")
		return errors.New("postID required")
	}
	return post.model.AddImpression(post.ID)
}

func (post *Post) Sync() error {
	if post.ID == 0 {
		return fmt.Errorf("Post.Sync(): required postID not set in post controller")
	}

	postData, err := post.model.GetByID(post.ID)
	if err != nil {
		return err
	}
	post.setFromModel(postData)
	return nil
}

func NewPostController(db *sql.DB) *Post {
	commentModel := model.NewCommentModel(db)
	commentController := &Comment{model: commentModel}

	return &Post{
		model:      model.NewPostModel(db),
		postAction: model.NewPostActionModel(db),
		comment:    commentController,
	}
}
