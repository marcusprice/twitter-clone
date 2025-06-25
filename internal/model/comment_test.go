package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/testhelpers"
	"github.com/marcusprice/twitter-clone/internal/testutil"
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
