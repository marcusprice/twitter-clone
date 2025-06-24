package controller

import (
	"database/sql"
	"strings"
	"time"

	"github.com/marcusprice/twitter-clone/internal/client"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/logger"
	"github.com/marcusprice/twitter-clone/internal/model"
	"github.com/marcusprice/twitter-clone/internal/util"
)

type DepthLimitError struct{}

func (d DepthLimitError) Error() string {
	return "Reply depth exceeds limit"
}

const DEPTH_LIMIT = 1

type Comment struct {
	model           *model.CommentModel
	post            *Post
	replyGuy        *client.ReplyGuyClient
	ID              int
	PostID          int
	Depth           int
	ParentCommentID int
	Content         string
	LikeCount       int
	RetweetCount    int
	BookmarkCount   int
	Impressions     int
	Image           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	Author          struct {
		Username    string
		DisplayName string
		Avatar      string
	}
	IsRetweet            bool
	RetweeterUsername    string
	RetweeterDisplayName string
}

func (comment *Comment) ByID(commentID int) (*Comment, error) {
	queriedComment := &Comment{}
	commentData, err := comment.model.GetByID(commentID)
	if err != nil {
		return &Comment{}, err
	}
	queriedComment.setFromModel(commentData)
	return queriedComment, nil
}

func (comment *Comment) setFromModel(commentData dtypes.CommentData) {
	comment.ID = commentData.ID
	comment.PostID = commentData.PostID
	comment.ParentCommentID = commentData.ParentCommentID
	comment.Content = commentData.Content
	comment.LikeCount = commentData.LikeCount
	comment.RetweetCount = commentData.RetweetCount
	comment.BookmarkCount = commentData.BookmarkCount
	comment.Impressions = commentData.Impressions
	comment.Image = commentData.Image
	comment.CreatedAt = util.ParseTime(commentData.CreatedAt)
	comment.UpdatedAt = util.ParseTime(commentData.UpdatedAt)
	comment.Author.Username = commentData.Author.Username
	comment.Author.DisplayName = commentData.Author.DisplayName
	comment.Author.Avatar = commentData.Author.Avatar
	comment.Depth = commentData.Depth
}

func (comment *Comment) New(commentInput dtypes.CommentInput) (*Comment, error) {
	var commentID int
	var err error
	if commentInput.ParentCommentID == 0 {
		commentID, err = comment.model.NewPostComment(commentInput)
	} else {
		parentComment, err := comment.model.GetByID(commentInput.ParentCommentID)
		if err != nil {
			return &Comment{}, err
		}

		if parentComment.Depth >= DEPTH_LIMIT {
			logger.LogWarn("Comment.New(): Reply depth exceeds limit")
			return &Comment{}, DepthLimitError{}
		}

		commentID, err = comment.model.NewCommentReply(commentInput)
	}

	if err != nil {
		return &Comment{}, err
	}

	commentData, err := comment.model.GetByID(commentID)
	if err != nil {
		return &Comment{}, err
	}

	newComment := &Comment{}
	newComment.setFromModel(commentData)

	for _, guy := range comment.replyGuy.GetReplyGuys() {
		if strings.Contains(newComment.Content, guy) {
			parentPost := comment.post
			err := parentPost.ByID(newComment.PostID)
			if err != nil {
				logger.LogError("Comment.New() error querying newComment.PostID: " + err.Error())
				break
			}

			parentComment := &Comment{}
			if newComment.ParentCommentID != 0 {
				// TODO: duplicate fetch. in refactor, reuse earlier query of newComment.ParentCommentID
				parentComment, err = comment.ByID(newComment.ParentCommentID)
				if err != nil {
					logger.LogError("Comment.New() error querying newComment.ParentCommentID: " + err.Error())
					break
				}
			}

			replyGuyRequest := dtypes.ReplyGuyRequest{
				PostID:                      parentPost.ID,
				PostAuthorUsername:          parentPost.Author.Username,
				PostContent:                 parentPost.Content,
				ParentCommentID:             parentComment.ID,
				ParentCommentAuthorUsername: parentComment.Author.Username,
				ParentCommentContent:        parentComment.Content,
				RequesterUsername:           newComment.Author.Username,
				Model:                       strings.TrimPrefix(guy, "@"),
				Prompt:                      newComment.Content,
			}
			err = comment.replyGuy.RequestReply(replyGuyRequest)
			if err != nil {
				logger.LogError("Comment.New() error with reply guy request: " + err.Error())
			}
		}
	}

	return newComment, nil
}

func NewCommentController(db *sql.DB) *Comment {
	replyGuy := client.NewReplyGuyClient()
	postController := NewPostController(db)
	model := model.NewCommentModel(db)
	return &Comment{model: model, post: postController, replyGuy: replyGuy}
}
