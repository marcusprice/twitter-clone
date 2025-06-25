package controller

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/model"
	"github.com/marcusprice/twitter-clone/internal/testhelpers"
	"github.com/marcusprice/twitter-clone/internal/testutil"
	"github.com/marcusprice/twitter-clone/internal/util"
)

func TestCommentByID(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, _ time.Time) {
		tu := testutil.NewTestUtil(t)
		commentInput := dtypes.CommentInput{
			PostID:  1,
			UserID:  1,
			Content: "Esters is the besters",
		}
		commentID := testhelpers.CreateComment(commentInput, db)
		commentModel := model.NewCommentModel(db)
		Comment := &Comment{model: commentModel, replyGuy: &testhelpers.MockReplyGuyClient{}}

		esteComment, err := Comment.ByID(commentID)
		tu.AssertErrorNil(err)
		tu.AssertEqual(commentInput.Content, esteComment.Content)
		tu.AssertEqual(commentInput.PostID, esteComment.PostID)
		tu.AssertEqual(commentInput.UserID, esteComment.UserID)

		notFoundComment, err := Comment.ByID(42069)
		tu.AssertErrorNotNil(err)
		tu.AssertEqual(0, notFoundComment.ID)
	})
}

func TestCommentNew(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, _ time.Time) {
		tu := testutil.NewTestUtil(t)
		model := model.NewCommentModel(db)
		Comment := &Comment{model: model, replyGuy: &testhelpers.MockReplyGuyClient{}}
		commentInput := dtypes.CommentInput{
			PostID:  1,
			UserID:  1,
			Content: "Freeskate broski",
		}
		newComment, err := Comment.New(commentInput)
		queriedComment := testhelpers.QueryComment(newComment.ID, db)
		tu.AssertErrorNil(err)
		tu.AssertEqual(commentInput.Content, newComment.Content)
		tu.AssertEqual(commentInput.PostID, newComment.PostID)
		tu.AssertEqual(commentInput.UserID, newComment.UserID)
		tu.AssertEqual(queriedComment.ID, newComment.ID)
		tu.AssertEqual(queriedComment.Content, newComment.Content)
		tu.AssertEqual(queriedComment.PostID, newComment.PostID)
		tu.AssertEqual(queriedComment.UserID, newComment.UserID)
		tu.AssertEqual(queriedComment.Author.Username, newComment.Author.Username)
		tu.AssertEqual(queriedComment.Author.DisplayName, newComment.Author.DisplayName)
		tu.AssertEqual(queriedComment.Author.Avatar, newComment.Author.Avatar)
		tu.AssertEqual(queriedComment.Retweeter.Username, newComment.RetweeterUsername)
		tu.AssertEqual(queriedComment.Retweeter.DisplayName, newComment.RetweeterDisplayName)
		tu.AssertEqual(queriedComment.Depth, newComment.Depth)
		tu.AssertEqual(queriedComment.ParentCommentID, newComment.ParentCommentID)
		tu.AssertEqual(queriedComment.LikeCount, newComment.LikeCount)
		tu.AssertEqual(queriedComment.RetweetCount, newComment.RetweetCount)
		tu.AssertEqual(queriedComment.BookmarkCount, newComment.BookmarkCount)
		tu.AssertEqual(queriedComment.Impressions, newComment.Impressions)
		tu.AssertEqual(queriedComment.Image, newComment.Image)
		tu.AssertEqual(util.ParseTime(queriedComment.CreatedAt), newComment.CreatedAt)
		tu.AssertEqual(util.ParseTime(queriedComment.UpdatedAt), newComment.UpdatedAt)
	})
}

func TestCommentNewReply(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, _ time.Time) {
		tu := testutil.NewTestUtil(t)
		model := model.NewCommentModel(db)
		Comment := &Comment{model: model, replyGuy: &testhelpers.MockReplyGuyClient{}}
		commentInput := dtypes.CommentInput{
			PostID:  1,
			UserID:  1,
			Content: "Freeskate broski",
		}
		topLevelComment, err := Comment.New(commentInput)
		commentInput = dtypes.CommentInput{
			PostID:          1,
			ParentCommentID: topLevelComment.ID,
			UserID:          2,
			Content:         "Did you call me broski?",
		}
		newComment, err := Comment.New(commentInput)
		queriedComment := testhelpers.QueryComment(newComment.ID, db)
		tu.AssertErrorNil(err)
		tu.AssertEqual(commentInput.Content, newComment.Content)
		tu.AssertEqual(commentInput.PostID, newComment.PostID)
		tu.AssertEqual(commentInput.UserID, newComment.UserID)
		tu.AssertEqual(topLevelComment.ID, newComment.ParentCommentID)
		tu.AssertEqual(queriedComment.ID, newComment.ID)
		tu.AssertEqual(queriedComment.Content, newComment.Content)
		tu.AssertEqual(queriedComment.PostID, newComment.PostID)
		tu.AssertEqual(queriedComment.UserID, newComment.UserID)
		tu.AssertEqual(queriedComment.Author.Username, newComment.Author.Username)
		tu.AssertEqual(queriedComment.Author.DisplayName, newComment.Author.DisplayName)
		tu.AssertEqual(queriedComment.Author.Avatar, newComment.Author.Avatar)
		tu.AssertEqual(queriedComment.Retweeter.Username, newComment.RetweeterUsername)
		tu.AssertEqual(queriedComment.Retweeter.DisplayName, newComment.RetweeterDisplayName)
		tu.AssertEqual(queriedComment.Depth, newComment.Depth)
		tu.AssertEqual(queriedComment.LikeCount, newComment.LikeCount)
		tu.AssertEqual(queriedComment.RetweetCount, newComment.RetweetCount)
		tu.AssertEqual(queriedComment.BookmarkCount, newComment.BookmarkCount)
		tu.AssertEqual(queriedComment.Impressions, newComment.Impressions)
		tu.AssertEqual(queriedComment.Image, newComment.Image)
		tu.AssertEqual(util.ParseTime(queriedComment.CreatedAt), newComment.CreatedAt)
		tu.AssertEqual(util.ParseTime(queriedComment.UpdatedAt), newComment.UpdatedAt)

		commentInput = dtypes.CommentInput{
			PostID:          1,
			ParentCommentID: newComment.ID,
			UserID:          1,
			Content:         "I sure did bucko",
		}
		badComment, err := Comment.New(commentInput)
		var depthLimitError DepthLimitError
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(errors.As(err, &depthLimitError))
		tu.AssertEqual("Reply depth exceeds limit", depthLimitError.Error())
		tu.AssertEqual(0, badComment.ID)
	})
}

func TestNewCommentWithReplyGuyRequest(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, _ time.Time) {
		tu := testutil.NewTestUtil(t)
		model := model.NewCommentModel(db)
		postController := NewPostController(db)
		replyGuyMockClient := &testhelpers.MockReplyGuyClient{}
		Comment := &Comment{
			model:    model,
			replyGuy: replyGuyMockClient,
			post:     postController,
		}
		donnaHayward := testhelpers.QueryUser(6, db)
		commentInput := dtypes.CommentInput{
			UserID:  donnaHayward.ID,
			PostID:  41,
			Content: "@dalecooper already questioned James, if that's what you're implying",
		}

		op := NewPostController(db)
		op.ByID(41)

		// TODO this ReplyGuyRequest struct sucks
		newComment, err := Comment.New(commentInput)
		calledWith := replyGuyMockClient.CalledWith
		tu.AssertErrorNil(err)
		tu.AssertEqual(newComment.PostID, calledWith.PostID)
		tu.AssertEqual(newComment.Author.Username, calledWith.RequesterUsername)
		tu.AssertEqual(op.Content, calledWith.PostContent)
		tu.AssertEqual(newComment.Content, calledWith.Prompt)
		tu.AssertEqual(0, calledWith.ParentCommentID)
		tu.AssertEqual("", calledWith.ParentCommentAuthorUsername)
		tu.AssertEqual("", calledWith.ParentCommentContent)
		tu.AssertEqual(donnaHayward.Username, calledWith.RequesterUsername)
		tu.AssertEqual("dalecooper", calledWith.Model)

		audreyHorne := testhelpers.QueryUser(4, db)
		commentInput = dtypes.CommentInput{
			UserID:          audreyHorne.ID,
			PostID:          41,
			Content:         "@dalecooper is this true?",
			ParentCommentID: newComment.ID,
		}
		commentReply, err := Comment.New(commentInput)
		calledWith = replyGuyMockClient.CalledWith
		tu.AssertEqual(op.Content, calledWith.PostContent)
		tu.AssertEqual(op.ID, calledWith.PostID)
		tu.AssertEqual(commentReply.PostID, calledWith.PostID)
		tu.AssertEqual(commentReply.Author.Username, calledWith.RequesterUsername)
		tu.AssertEqual(commentReply.Content, calledWith.Prompt)
		tu.AssertEqual(commentReply.ParentCommentID, calledWith.ParentCommentID)
		tu.AssertEqual(donnaHayward.Username, calledWith.ParentCommentAuthorUsername)
		tu.AssertEqual(newComment.Content, calledWith.ParentCommentContent)
		tu.AssertEqual(audreyHorne.Username, calledWith.RequesterUsername)
		tu.AssertEqual("dalecooper", calledWith.Model)
	})
}
