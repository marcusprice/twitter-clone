package api

import (
	"time"

	"github.com/marcusprice/twitter-clone/internal/controller"
)

type TimelinePayload struct {
	Posts          []PostPayload `json:"posts"`
	HasMore        bool          `json:"hasMore"`
	PostsRemaining int           `json:"postsRemaining"`
}

type UserPayload struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	DisplayName string `json:"displayName"`
}

type PostPayload struct {
	ID            int           `json:"postID"`
	Content       string        `json:"content"`
	LikeCount     int           `json:"likeCount"`
	RetweetCount  int           `json:"retweetCount"`
	BookmarkCount int           `json:"bookmarkCount"`
	Impressions   int           `json:"impressions"`
	Image         string        `json:"image"`
	CreatedAt     time.Time     `json:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt"`
	Author        AuthorPayload `json:"author"`
}

type AuthorPayload struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Avatar      string `json:"avatar"`
}

func generatePostPayload(post *controller.Post) PostPayload {
	author := AuthorPayload{
		Username:    post.Author.Username,
		DisplayName: post.Author.DisplayName,
		Avatar:      post.Author.Avatar,
	}

	return PostPayload{
		ID:            post.ID,
		Content:       post.Content,
		LikeCount:     post.LikeCount,
		RetweetCount:  post.RetweetCount,
		BookmarkCount: post.BookmarkCount,
		Impressions:   post.Impressions,
		Image:         post.Image,
		CreatedAt:     post.CreatedAt,
		UpdatedAt:     post.UpdatedAt,
		Author:        author,
	}
}
