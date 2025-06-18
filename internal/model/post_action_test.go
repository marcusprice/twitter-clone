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

func TestPostActionUnlike(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		postAction := NewPostActionModel(db)

		insertPostLikeRow(1, 1, db, t)
		insertPostLikeRow(1, 2, db, t)
		insertPostLikeRow(1, 3, db, t)
		insertPostLikeRow(1, 4, db, t)

		err := postAction.Unlike(1, 4)
		tu.AssertErrorNil(err)
		postData := queryPost(1, db)
		tu.AssertEqual(3, postData.LikeCount)

		err = postAction.Unlike(1, 3)
		tu.AssertErrorNil(err)
		postData = queryPost(1, db)
		tu.AssertEqual(2, postData.LikeCount)

		err = postAction.Unlike(1, 2)
		tu.AssertErrorNil(err)
		postData = queryPost(1, db)
		tu.AssertEqual(1, postData.LikeCount)

		err = postAction.Unlike(1, 2)
		tu.AssertErrorNil(err)
		postData = queryPost(1, db)
		tu.AssertEqual(1, postData.LikeCount)

		err = postAction.Unlike(1, 1)
		tu.AssertErrorNil(err)
		postData = queryPost(1, db)
		tu.AssertEqual(0, postData.LikeCount)

		err = postAction.Unlike(1, 1)
		tu.AssertErrorNil(err)
		postData = queryPost(1, db)
		tu.AssertEqual(0, postData.LikeCount)
	})
}

func TestPostActionLikeConstraintError(t *testing.T) {
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

func TestPostActionRetweet(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		postAction := NewPostActionModel(db)
		userID := 3
		postID := 2
		err := postAction.Retweet(postID, userID)
		tu.AssertErrorNil(err)

		queriedPostID, queriedUserID := queryPostRetweetRowByID(1, db)
		tu.AssertEqual(postID, queriedPostID)
		tu.AssertEqual(userID, queriedUserID)
	})
}

func TestPostActionUnRetweet(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		postAction := NewPostActionModel(db)

		insertPostRetweetRow(1, 1, db, t)
		insertPostRetweetRow(1, 2, db, t)
		insertPostRetweetRow(1, 3, db, t)
		insertPostRetweetRow(1, 4, db, t)

		err := postAction.UnRetweet(1, 4)
		tu.AssertErrorNil(err)
		postData := queryPost(1, db)
		tu.AssertEqual(3, postData.RetweetCount)

		err = postAction.UnRetweet(1, 3)
		tu.AssertErrorNil(err)
		postData = queryPost(1, db)
		tu.AssertEqual(2, postData.RetweetCount)

		err = postAction.UnRetweet(1, 2)
		tu.AssertErrorNil(err)
		postData = queryPost(1, db)
		tu.AssertEqual(1, postData.RetweetCount)

		err = postAction.UnRetweet(1, 2)
		tu.AssertErrorNil(err)
		postData = queryPost(1, db)
		tu.AssertEqual(1, postData.RetweetCount)

		err = postAction.UnRetweet(1, 1)
		tu.AssertErrorNil(err)
		postData = queryPost(1, db)
		tu.AssertEqual(0, postData.RetweetCount)

		err = postAction.UnRetweet(1, 1)
		tu.AssertErrorNil(err)
		postData = queryPost(1, db)
		tu.AssertEqual(0, postData.RetweetCount)
	})
}

func TestPostActionRetweetConstraintError(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		postAction := &PostAction{db}
		var constraintError dbutils.ConstraintError

		err := postAction.Retweet(1, 42069)
		postData := queryPost(1, db)
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(errors.As(err, &constraintError))
		tu.AssertEqual(dbutils.FOREIGN_KEY_ERROR, constraintError.Constraint)
		tu.AssertEqual(0, postData.RetweetCount)

		err = postAction.Retweet(42069, 1)
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(errors.As(err, &constraintError))
		tu.AssertEqual(dbutils.FOREIGN_KEY_ERROR, constraintError.Constraint)

		err = postAction.Retweet(42069, 42069)
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(errors.As(err, &constraintError))
		tu.AssertEqual(dbutils.FOREIGN_KEY_ERROR, constraintError.Constraint)
	})
}

func TestPostActionBookmark(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		postAction := NewPostActionModel(db)
		userID := 3
		postID := 2
		err := postAction.Bookmark(postID, userID)
		tu.AssertErrorNil(err)

		queriedPostID, queriedUserID := queryPostBookmarkRowByID(1, db)
		tu.AssertEqual(postID, queriedPostID)
		tu.AssertEqual(userID, queriedUserID)
	})
}

func TestPostActionUnBookmark(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		postAction := NewPostActionModel(db)

		insertPostBookmarkRow(1, 1, db, t)
		insertPostBookmarkRow(1, 2, db, t)
		insertPostBookmarkRow(1, 3, db, t)
		insertPostBookmarkRow(1, 4, db, t)

		err := postAction.UnBookmark(1, 4)
		tu.AssertErrorNil(err)
		postData := queryPost(1, db)
		tu.AssertEqual(3, postData.BookmarkCount)

		err = postAction.UnBookmark(1, 3)
		tu.AssertErrorNil(err)
		postData = queryPost(1, db)
		tu.AssertEqual(2, postData.BookmarkCount)

		err = postAction.UnBookmark(1, 2)
		tu.AssertErrorNil(err)
		postData = queryPost(1, db)
		tu.AssertEqual(1, postData.BookmarkCount)

		err = postAction.UnBookmark(1, 2)
		tu.AssertErrorNil(err)
		postData = queryPost(1, db)
		tu.AssertEqual(1, postData.BookmarkCount)

		err = postAction.UnBookmark(1, 1)
		tu.AssertErrorNil(err)
		postData = queryPost(1, db)
		tu.AssertEqual(0, postData.BookmarkCount)

		err = postAction.UnBookmark(1, 1)
		tu.AssertErrorNil(err)
		postData = queryPost(1, db)
		tu.AssertEqual(0, postData.BookmarkCount)
	})
}

func TestPostActionBookmarkConstraintError(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		postAction := &PostAction{db}
		var constraintError dbutils.ConstraintError

		err := postAction.Bookmark(1, 42069)
		postData := queryPost(1, db)
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(errors.As(err, &constraintError))
		tu.AssertEqual(dbutils.FOREIGN_KEY_ERROR, constraintError.Constraint)
		tu.AssertEqual(0, postData.BookmarkCount)

		err = postAction.Bookmark(42069, 1)
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(errors.As(err, &constraintError))
		tu.AssertEqual(dbutils.FOREIGN_KEY_ERROR, constraintError.Constraint)

		err = postAction.Bookmark(42069, 42069)
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

func queryPostRetweetRowByID(id int, db *sql.DB) (int, int) {
	var postID int
	var userID int

	db.
		QueryRow("SELECT post_id, user_id FROM PostRetweet WHERE id = $1;", id).
		Scan(&postID, &userID)

	return postID, userID
}

func queryPostBookmarkRowByID(id int, db *sql.DB) (int, int) {
	var postID int
	var userID int

	db.
		QueryRow("SELECT post_id, user_id FROM PostBookmark WHERE id = $1;", id).
		Scan(&postID, &userID)

	return postID, userID
}

func insertPostLikeRow(postID, userID int, db *sql.DB, t *testing.T) {
	_, err := db.Exec(`
		INSERT INTO PostLike 
			(post_id, user_id)
		VALUES
			($1, $2);
	`, postID, userID)

	if err != nil {
		t.Fatal("Error inserting PostLike row", err.Error())
	}
}

func insertPostRetweetRow(postID, userID int, db *sql.DB, t *testing.T) {
	_, err := db.Exec(`
		INSERT INTO PostRetweet 
			(post_id, user_id)
		VALUES
			($1, $2);
	`, postID, userID)

	if err != nil {
		t.Fatal("Error inserting PostRetweet row", err.Error())
	}
}

func insertPostBookmarkRow(postID, userID int, db *sql.DB, t *testing.T) {
	_, err := db.Exec(`
		INSERT INTO PostBookmark
			(post_id, user_id)
		VALUES
			($1, $2);
	`, postID, userID)

	if err != nil {
		t.Fatal("Error inserting PostRetweet row", err.Error())
	}
}
