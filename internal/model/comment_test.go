package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/dbutils"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/testhelpers"
	"github.com/marcusprice/twitter-clone/internal/testutil"
	"github.com/marcusprice/twitter-clone/internal/util"
)

func TestCommentGetByID(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, _ time.Time) {
		tu := testutil.NewTestUtil(t)
		CommentModel := &CommentModel{db: db}
		commentInput := dtypes.CommentInput{
			PostID:  1,
			UserID:  1,
			Content: "Damn now I want to buy chocolate",
		}
		testhelpers.CreateComment(commentInput, db)

		commentData, err := CommentModel.GetByID(1)
		expected, _ := json.Marshal(testhelpers.QueryComment(1, db))
		actual, _ := json.Marshal(commentData)

		tu.AssertErrorNil(err)
		tu.AssertEqual(1, commentData.ID)
		tu.AssertEqual(commentInput.PostID, commentData.PostID)
		tu.AssertEqual(commentInput.UserID, commentData.UserID)
		tu.AssertEqual(commentInput.Content, commentData.Content)
		tu.AssertEqual(commentInput.Image, commentData.Image)
		tu.AssertEqual(string(expected), string(actual))

		commentData, err = CommentModel.GetByID(42069)
		var commentNotFoundError CommentNotFoundError
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(errors.As(err, &commentNotFoundError))
		tu.AssertEqual("Comment not found", commentNotFoundError.Error())
		tu.AssertEqual(0, commentData.ID)
	})
}

func TestNewPostComment(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, _ time.Time) {
		tu := testutil.NewTestUtil(t)
		CommentModel := &CommentModel{db: db}

		wallphace := testhelpers.QueryUser(1, db)
		commentInput := dtypes.CommentInput{
			PostID:  1,
			UserID:  wallphace.ID,
			Content: "Damn now I want to buy chocolate",
			Image:   "suhdude.jpeg",
		}

		beforeInsert := time.Now().UTC().Add(-1 * time.Minute)
		rowID, err := CommentModel.NewPostComment(commentInput)
		afterInsert := time.Now().UTC().Add(time.Minute)
		commentData := testhelpers.QueryComment(rowID, db)
		tu.AssertErrorNil(err)
		tu.AssertEqual(1, rowID)
		tu.AssertEqual(rowID, commentData.ID)
		tu.AssertEqual(commentInput.Content, commentData.Content)
		tu.AssertEqual(commentInput.Image, commentData.Image)
		tu.AssertEqual(commentInput.UserID, commentData.UserID)
		tu.AssertEqual(commentInput.PostID, commentData.PostID)
		tu.AssertEqual(0, commentData.ParentCommentID)
		tu.AssertEqual(0, commentData.Depth)
		tu.AssertEqual(0, commentData.LikeCount)
		tu.AssertEqual(0, commentData.RetweetCount)
		tu.AssertEqual(0, commentData.BookmarkCount)
		tu.AssertEqual(0, commentData.Impressions)
		tu.AssertEqual(wallphace.Username, commentData.Author.Username)
		tu.AssertEqual(wallphace.DisplayName, commentData.Author.DisplayName)
		tu.AssertEqual(wallphace.Avatar, commentData.Author.Avatar)
		tu.AssertTrue(beforeInsert.Before(util.ParseTime(commentData.CreatedAt)))
		tu.AssertTrue(afterInsert.After(util.ParseTime(commentData.CreatedAt)))
		tu.AssertTrue(beforeInsert.Before(util.ParseTime(commentData.UpdatedAt)))
		tu.AssertTrue(afterInsert.After(util.ParseTime(commentData.UpdatedAt)))

		commentInput = dtypes.CommentInput{
			PostID:  42069,
			UserID:  wallphace.ID,
			Content: "Damn now I want to buy chocolate",
			Image:   "suhdude.jpeg",
		}
		rowID, err = CommentModel.NewPostComment(commentInput)
		var constraintError dbutils.ConstraintError
		tu.AssertErrorNotNil(err)
		tu.AssertEqual(-1, rowID)
		tu.AssertTrue(errors.As(err, &constraintError))
		tu.AssertEqual(dbutils.FOREIGN_KEY_ERROR, constraintError.Constraint)

		commentInput = dtypes.CommentInput{
			PostID:  1,
			UserID:  42069,
			Content: "Damn now I want to buy chocolate",
			Image:   "suhdude.jpeg",
		}
		rowID, err = CommentModel.NewPostComment(commentInput)
		tu.AssertErrorNotNil(err)
		tu.AssertEqual(-1, rowID)
		tu.AssertTrue(errors.As(err, &constraintError))
		tu.AssertEqual(dbutils.FOREIGN_KEY_ERROR, constraintError.Constraint)

		commentInput = dtypes.CommentInput{
			PostID:  1,
			UserID:  1,
			Content: "",
			Image:   "",
		}
		rowID, err = CommentModel.NewPostComment(commentInput)
		tu.AssertErrorNotNil(err)
		tu.AssertEqual(-1, rowID)
		tu.AssertTrue(errors.As(err, &constraintError))
		tu.AssertEqual(dbutils.CHECK_ERROR, constraintError.Constraint)
	})
}

func TestNewCommentReply(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, _ time.Time) {
		tu := testutil.NewTestUtil(t)
		CommentModel := &CommentModel{db: db}

		wallphace := testhelpers.QueryUser(1, db)
		parentComment := dtypes.CommentInput{
			PostID:  1,
			UserID:  wallphace.ID,
			Content: "Damn now I want to buy chocolate",
			Image:   "suhdude.jpeg",
		}
		parentCommentID := testhelpers.CreateComment(parentComment, db)

		commentReplyInput := dtypes.CommentInput{
			PostID:          1,
			ParentCommentID: parentCommentID,
			UserID:          wallphace.ID,
			Content:         "Damn now I want to buy chocolate",
			Image:           "suhdude.jpeg",
		}

		beforeInsert := time.Now().UTC().Add(-1 * time.Minute)
		rowID, err := CommentModel.NewCommentReply(commentReplyInput)
		afterInsert := time.Now().UTC().Add(time.Minute)
		commentData := testhelpers.QueryComment(rowID, db)
		tu.AssertErrorNil(err)
		tu.AssertEqual(2, rowID)
		tu.AssertEqual(rowID, commentData.ID)
		tu.AssertEqual(commentReplyInput.Content, commentData.Content)
		tu.AssertEqual(commentReplyInput.Image, commentData.Image)
		tu.AssertEqual(commentReplyInput.UserID, commentData.UserID)
		tu.AssertEqual(commentReplyInput.PostID, commentData.PostID)
		tu.AssertEqual(1, commentData.ParentCommentID)
		tu.AssertEqual(1, commentData.Depth)
		tu.AssertEqual(0, commentData.LikeCount)
		tu.AssertEqual(0, commentData.RetweetCount)
		tu.AssertEqual(0, commentData.BookmarkCount)
		tu.AssertEqual(0, commentData.Impressions)
		tu.AssertEqual(wallphace.Username, commentData.Author.Username)
		tu.AssertEqual(wallphace.DisplayName, commentData.Author.DisplayName)
		tu.AssertEqual(wallphace.Avatar, commentData.Author.Avatar)
		tu.AssertTrue(beforeInsert.Before(util.ParseTime(commentData.CreatedAt)))
		tu.AssertTrue(afterInsert.After(util.ParseTime(commentData.CreatedAt)))
		tu.AssertTrue(beforeInsert.Before(util.ParseTime(commentData.UpdatedAt)))
		tu.AssertTrue(afterInsert.After(util.ParseTime(commentData.UpdatedAt)))

		commentReplyInput = dtypes.CommentInput{
			PostID:          42069,
			UserID:          wallphace.ID,
			ParentCommentID: parentCommentID,
			Content:         "Damn now I want to buy chocolate",
			Image:           "suhdude.jpeg",
		}
		rowID, err = CommentModel.NewCommentReply(commentReplyInput)
		var constraintError dbutils.ConstraintError
		tu.AssertErrorNotNil(err)
		tu.AssertEqual(-1, rowID)
		tu.AssertTrue(errors.As(err, &constraintError))
		tu.AssertEqual(dbutils.FOREIGN_KEY_ERROR, constraintError.Constraint)

		commentReplyInput = dtypes.CommentInput{
			PostID:          1,
			UserID:          42069,
			ParentCommentID: parentCommentID,
			Content:         "Damn now I want to buy chocolate",
			Image:           "suhdude.jpeg",
		}
		rowID, err = CommentModel.NewCommentReply(commentReplyInput)
		tu.AssertErrorNotNil(err)
		tu.AssertEqual(-1, rowID)
		tu.AssertTrue(errors.As(err, &constraintError))
		tu.AssertEqual(dbutils.FOREIGN_KEY_ERROR, constraintError.Constraint)

		commentReplyInput = dtypes.CommentInput{
			PostID:          1,
			UserID:          1,
			ParentCommentID: 42069,
			Content:         "Damn now I want to buy chocolate",
			Image:           "suhdude.jpeg",
		}
		rowID, err = CommentModel.NewCommentReply(commentReplyInput)
		tu.AssertErrorNotNil(err)
		tu.AssertEqual(-1, rowID)
		tu.AssertTrue(errors.As(err, &constraintError))
		tu.AssertEqual(dbutils.FOREIGN_KEY_ERROR, constraintError.Constraint)

		commentReplyInput = dtypes.CommentInput{
			PostID:          1,
			UserID:          1,
			ParentCommentID: parentCommentID,
			Content:         "",
			Image:           "",
		}
		rowID, err = CommentModel.NewCommentReply(commentReplyInput)
		tu.AssertErrorNotNil(err)
		tu.AssertEqual(-1, rowID)
		tu.AssertTrue(errors.As(err, &constraintError))
		tu.AssertEqual(dbutils.CHECK_ERROR, constraintError.Constraint)
	})
}
