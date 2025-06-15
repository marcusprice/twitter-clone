package controller

import (
	"database/sql"
	"time"

	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/model"
	"github.com/marcusprice/twitter-clone/internal/util"
)

type Post struct {
	model         *model.PostModel
	ID            int
	UserID        int
	Content       string
	LikeCount     int
	RetweetCount  int
	BookmarkCount int
	Impressions   int
	Image         string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Author        struct {
		Username    string
		DisplayName string
		Avatar      string
	}
}

func (p *Post) setFromModel(postData model.PostData) {
	p.ID = postData.ID
	p.UserID = postData.UserID
	p.Content = postData.Content
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

func NewPostController(db *sql.DB) *Post {
	return &Post{
		model: model.NewPostModel(db),
	}
}
