package model

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/dbutils"
	"github.com/marcusprice/twitter-clone/internal/testutil"
)

func TestPostActionLike(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		postAction := NewPostActionModel(db)
		userID := 3
		postID := 2
		err := postAction.Like(postID, userID)
		tu.AssertErrorNil(err)

		queriedPostID, queriedUserID := queryPostLikeRowByID(1, db)
		tu.AssertEqual(postID, queriedPostID)
		tu.AssertEqual(userID, queriedUserID)
	})
}

func TestPostActionConstraintError(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		postAction := &PostAction{db}
		var constraintError dbutils.ConstraintError

		err := postAction.Like(1, 42069)
		postData := queryPost(1, db)
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(errors.As(err, &constraintError))
		tu.AssertEqual(dbutils.FOREIGN_KEY_ERROR, constraintError.Constraint)
		tu.AssertEqual(0, postData.LikeCount)

		err = postAction.Like(42069, 1)
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(errors.As(err, &constraintError))
		tu.AssertEqual(dbutils.FOREIGN_KEY_ERROR, constraintError.Constraint)

		err = postAction.Like(42069, 42069)
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(errors.As(err, &constraintError))
		tu.AssertEqual(dbutils.FOREIGN_KEY_ERROR, constraintError.Constraint)
	})
}

func queryPostLikeRowByID(id int, db *sql.DB) (int, int) {
	var postID int
	var userID int

	db.
		QueryRow("SELECT post_id, user_id FROM PostLike WHERE id = $1;", id).
		Scan(&postID, &userID)

	return postID, userID
}
