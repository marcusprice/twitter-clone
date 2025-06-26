package api

import (
	"fmt"
	"os"
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

type AuthorPayload struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Avatar      string `json:"avatar"`
}

type PostPayload struct {
	ID                   int           `json:"postID"`
	Content              string        `json:"content"`
	LikeCount            int           `json:"likeCount"`
	RetweetCount         int           `json:"retweetCount"`
	BookmarkCount        int           `json:"bookmarkCount"`
	Impressions          int           `json:"impressions"`
	Image                string        `json:"image"`
	CreatedAt            time.Time     `json:"createdAt"`
	UpdatedAt            time.Time     `json:"updatedAt"`
	Author               AuthorPayload `json:"author"`
	IsRetweet            bool          `json:"isRetweet"`
	RetweeterUsername    string        `json:"retweeterUsername"`
	RetweeterDisplayName string        `json:"retweeterDisplayName"`
}

func generatePostPayload(post *controller.Post) PostPayload {
	author := AuthorPayload{
		Username:    post.Author.Username,
		DisplayName: post.Author.DisplayName,
		Avatar:      post.Author.Avatar,
	}

	imageURL := fmt.Sprintf(
		"http://%s:%s%s%s",
		os.Getenv("HOST"),
		os.Getenv("PORT"),
		UPLOADS_PREFIX,
		post.Image)

	return PostPayload{
		ID:                   post.ID,
		Content:              post.Content,
		LikeCount:            post.LikeCount,
		RetweetCount:         post.RetweetCount,
		BookmarkCount:        post.BookmarkCount,
		Impressions:          post.Impressions,
		Image:                imageURL,
		CreatedAt:            post.CreatedAt,
		UpdatedAt:            post.UpdatedAt,
		Author:               author,
		IsRetweet:            post.Retweeter.Username != "",
		RetweeterUsername:    post.Retweeter.Username,
		RetweeterDisplayName: post.Retweeter.DisplayName,
	}
}

type CommentPayload struct {
	ID                   int           `json:"commentID"`
	PostID               int           `json:"postID"`
	ParentCommentID      int           `json:"parentCommentID"`
	Content              string        `json:"content"`
	LikeCount            int           `json:"likeCount"`
	RetweetCount         int           `json:"retweetCount"`
	BookmarkCount        int           `json:"bookmarkCount"`
	Impressions          int           `json:"impressions"`
	Image                string        `json:"image"`
	CreatedAt            time.Time     `json:"createdAt"`
	UpdatedAt            time.Time     `json:"updatedAt"`
	Author               AuthorPayload `json:"author"`
	IsRetweet            bool          `json:"isRetweet"`
	RetweeterUsername    string        `json:"retweeterUsername"`
	RetweeterDisplayName string        `json:"retweeterDisplayName"`
}

func generateCommentPayload(comment *controller.Comment) *CommentPayload {
	author := AuthorPayload{
		Username:    comment.Author.Username,
		DisplayName: comment.Author.DisplayName,
		Avatar:      comment.Author.Avatar,
	}

	return &CommentPayload{
		ID:                   comment.ID,
		PostID:               comment.PostID,
		ParentCommentID:      comment.ParentCommentID,
		Content:              comment.Content,
		LikeCount:            comment.LikeCount,
		RetweetCount:         comment.RetweetCount,
		BookmarkCount:        comment.BookmarkCount,
		Impressions:          comment.Impressions,
		Image:                comment.Image,
		CreatedAt:            comment.CreatedAt,
		UpdatedAt:            comment.UpdatedAt,
		IsRetweet:            comment.IsRetweet,
		RetweeterUsername:    comment.RetweeterUsername,
		RetweeterDisplayName: comment.RetweeterDisplayName,
		Author:               author,
	}
}

type CommentFromPostPayload struct {
	ID              int                       `json:"commentID"`
	PostID          int                       `json:"postID"`
	ParentCommentID int                       `json:"parentCommentID"`
	Content         string                    `json:"content"`
	LikeCount       int                       `json:"likeCount"`
	RetweetCount    int                       `json:"retweetCount"`
	BookmarkCount   int                       `json:"bookmarkCount"`
	Impressions     int                       `json:"impressions"`
	Image           string                    `json:"image"`
	CreatedAt       time.Time                 `json:"createdAt"`
	UpdatedAt       time.Time                 `json:"updatedAt"`
	Author          AuthorPayload             `json:"author"`
	Replies         []*CommentFromPostPayload `json:"replies"`
}

type PostAndCommentsPayload struct {
	ID            int                       `json:"postID"`
	Content       string                    `json:"content"`
	LikeCount     int                       `json:"likeCount"`
	RetweetCount  int                       `json:"retweetCount"`
	BookmarkCount int                       `json:"bookmarkCount"`
	Impressions   int                       `json:"impressions"`
	Image         string                    `json:"image"`
	CreatedAt     time.Time                 `json:"createdAt"`
	UpdatedAt     time.Time                 `json:"updatedAt"`
	Author        AuthorPayload             `json:"author"`
	Comments      []*CommentFromPostPayload `json:"comments"`
}

func generatePostAndCommentsPayload(post *controller.Post) PostAndCommentsPayload {
	postAndCommentsPayload := PostAndCommentsPayload{}
	postAndCommentsPayload.Comments = []*CommentFromPostPayload{}
	for _, comment := range post.Comments {
		commentPayload := &CommentFromPostPayload{}
		repliesPayload := []*CommentFromPostPayload{}
		for _, reply := range comment.Replies {
			replyPayload := &CommentFromPostPayload{}
			replyPayload.ID = reply.ID
			replyPayload.PostID = reply.PostID
			replyPayload.ParentCommentID = reply.ParentCommentID
			replyPayload.Content = reply.Content
			replyPayload.LikeCount = reply.LikeCount
			replyPayload.RetweetCount = reply.RetweetCount
			replyPayload.BookmarkCount = reply.BookmarkCount
			replyPayload.Impressions = reply.Impressions
			replyPayload.Image = reply.Image
			replyPayload.CreatedAt = reply.CreatedAt
			replyPayload.UpdatedAt = reply.UpdatedAt
			replyPayload.Author = AuthorPayload(reply.Author)
			repliesPayload = append(repliesPayload, replyPayload)
		}

		commentPayload.ID = comment.ID
		commentPayload.PostID = comment.PostID
		commentPayload.ParentCommentID = comment.ParentCommentID
		commentPayload.Content = comment.Content
		commentPayload.LikeCount = comment.LikeCount
		commentPayload.RetweetCount = comment.RetweetCount
		commentPayload.BookmarkCount = comment.BookmarkCount
		commentPayload.Impressions = comment.Impressions
		commentPayload.Image = comment.Image
		commentPayload.CreatedAt = comment.CreatedAt
		commentPayload.UpdatedAt = comment.UpdatedAt
		commentPayload.Author = AuthorPayload(comment.Author)
		commentPayload.Replies = repliesPayload

		postAndCommentsPayload.Comments = append(
			postAndCommentsPayload.Comments,
			commentPayload,
		)
	}

	postAndCommentsPayload.ID = post.ID
	postAndCommentsPayload.Content = post.Content
	postAndCommentsPayload.LikeCount = post.LikeCount
	postAndCommentsPayload.RetweetCount = post.RetweetCount
	postAndCommentsPayload.BookmarkCount = post.BookmarkCount
	postAndCommentsPayload.Impressions = post.Impressions
	postAndCommentsPayload.Image = post.Image
	postAndCommentsPayload.CreatedAt = post.CreatedAt
	postAndCommentsPayload.UpdatedAt = post.UpdatedAt
	postAndCommentsPayload.Author = AuthorPayload(post.Author)

	return postAndCommentsPayload
}
