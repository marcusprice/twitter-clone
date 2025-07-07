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
	model                *model.CommentModel
	post                 *Post
	replyGuy             client.ReplyGuyRequester
	ID                   int
	PostID               int
	UserID               int
	Depth                int
	ParentCommentID      int
	Content              string
	LikeCount            int
	RetweetCount         int
	BookmarkCount        int
	Impressions          int
	Image                string
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Author               dtypes.Author
	IsRetweet            bool
	RetweeterUsername    string
	RetweeterDisplayName string
	Replies              []*Comment
}

func (comment *Comment) setFromModel(commentData dtypes.CommentData) {
	comment.ID = commentData.ID
	comment.PostID = commentData.PostID
	comment.UserID = commentData.UserID
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

func (comment *Comment) ByID(commentID int) (*Comment, error) {
	queriedComment := &Comment{}
	commentData, err := comment.model.GetByID(commentID)
	if err != nil {
		return &Comment{}, err
	}
	queriedComment.setFromModel(commentData)
	return queriedComment, nil
}

func (comment *Comment) GetPostComments(postID int) ([]*Comment, error) {
	commentData, err := comment.model.GetByPostID(postID)
	if err != nil {
		return []*Comment{}, err
	}

	topLevelComments := []*Comment{}
	commentByParentMap := make(map[int][]*Comment)
	for _, c := range commentData {
		comment := &Comment{}
		comment.setFromModel(c)

		if comment.ParentCommentID != 0 {
			if slice, ok := commentByParentMap[comment.ParentCommentID]; !ok {
				commentByParentMap[comment.ParentCommentID] = []*Comment{comment}
			} else {
				slice = append(slice, comment)
				commentByParentMap[comment.ParentCommentID] = slice
			}

			continue
		}

		topLevelComments = append(topLevelComments, comment)
	}

	for _, comment := range topLevelComments {
		comment.Replies = []*Comment{}
		if replies, ok := commentByParentMap[comment.ID]; ok {
			comment.Replies = replies
		}
	}

	return topLevelComments, nil
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
			err := comment.handleReplyGuyRequest(guy, newComment)

			if err != nil {
				break
			}
		}
	}

	return newComment, nil
}

func (comment *Comment) handleReplyGuyRequest(guy string, newComment *Comment) error {
	parentPost := comment.post
	err := parentPost.ByID(newComment.PostID)
	if err != nil {
		logger.LogError("Comment.New() error querying newComment.PostID: " + err.Error())
		return err
	}

	parentComment := &Comment{}
	if newComment.ParentCommentID != 0 {
		// TODO: duplicate fetch. in refactor, reuse earlier query of newComment.ParentCommentID
		parentComment, err = comment.ByID(newComment.ParentCommentID)
		if err != nil {
			logger.LogError("Comment.New() error querying newComment.ParentCommentID: " + err.Error())
			return err
		}
	}

	replyGuyComment := dtypes.ReplyGuyComment{
		ID:      newComment.ID,
		Content: newComment.Content,
		Author:  newComment.Author,
	}

	replyGuyParentPostAuthor := dtypes.Author{
		Username: parentPost.Author.Username,
	}

	replyGuyParentPost := dtypes.ReplyGuyPost{
		ID:      parentPost.ID,
		Content: parentPost.Content,
		Author:  replyGuyParentPostAuthor,
	}

	var replyGuyParentComment dtypes.ReplyGuyComment
	var replyGuyParentCommentAuthor dtypes.Author
	if parentComment.ID != 0 {
		replyGuyParentCommentAuthor = dtypes.Author{
			Username: parentComment.Author.Username,
		}

		replyGuyParentComment = dtypes.ReplyGuyComment{
			ID:      parentComment.ID,
			Content: parentComment.Content,
			Author:  replyGuyParentCommentAuthor,
		}
	}

	replyGuyRequest := dtypes.ReplyGuyRequest{
		Model:         strings.TrimPrefix(guy, "@"),
		Comment:       replyGuyComment,
		ParentPost:    replyGuyParentPost,
		ParentComment: replyGuyParentComment,
	}

	if comment.replyGuy.RunAsync() {
		go comment.replyGuy.RequestReply(replyGuyRequest)
	} else {
		comment.replyGuy.RequestReply(replyGuyRequest)
	}

	return nil
}

func NewCommentController(db *sql.DB) *Comment {
	replyGuy := client.NewReplyGuyClient()
	postModel := model.NewPostModel(db)
	postController := &Post{model: postModel}

	model := model.NewCommentModel(db)
	return &Comment{model: model, post: postController, replyGuy: replyGuy}
}
