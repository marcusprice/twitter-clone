package model

import (
	"database/sql"
	"testing"
	"time"

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

func queryPostLikeRowByID(id int, db *sql.DB) (int, int) {
	var postID int
	var userID int

	db.
		QueryRow("SELECT post_id, user_id FROM PostLike WHERE id = $1;", id).
		Scan(&postID, &userID)

	return postID, userID
}
