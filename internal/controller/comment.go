package controller

import (
	"database/sql"
	"time"

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
	ID              int
	PostID          int
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

	return newComment, nil
}

func NewCommentController(db *sql.DB) *Comment {
	commentModel := model.NewCommentModel(db)
	return &Comment{model: commentModel}
}
