package controller

import (
	"database/sql"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/model"
	"github.com/marcusprice/twitter-clone/internal/testhelpers"
	"github.com/marcusprice/twitter-clone/internal/testutil"
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
		Comment := &Comment{model: commentModel}

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
