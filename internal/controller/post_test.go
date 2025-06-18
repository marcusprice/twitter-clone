package controller

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/constants"
	"github.com/marcusprice/twitter-clone/internal/dbutils"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/model"
	"github.com/marcusprice/twitter-clone/internal/testutil"
)

func TestSetFromModel(t *testing.T) {
	tu := testutil.NewTestUtil(t)
	post := &Post{}
	postAuthor := model.PostAuthor{
		Username:    "estecat",
		DisplayName: "Bubba",
		Avatar:      "lazy-cat.png",
	}
	postData := model.PostData{
		ID:            42069,
		UserID:        69,
		Content:       "Is it time for dinner yet?",
		Image:         "empty-dish.png",
		Author:        postAuthor,
		LikeCount:     1200,
		RetweetCount:  80,
		BookmarkCount: 12,
		Impressions:   10000,
		CreatedAt:     "2024-04-12 11:37:46",
		UpdatedAt:     "2024-05-18 08:02:13",
	}

	post.setFromModel(postData)
	tu.AssertEqual(42069, post.ID)
	tu.AssertEqual(69, post.UserID)
	tu.AssertEqual("estecat", post.Author.Username)
	tu.AssertEqual("Bubba", post.Author.DisplayName)
	tu.AssertEqual("lazy-cat.png", post.Author.Avatar)
	tu.AssertEqual("Is it time for dinner yet?", post.Content)
	tu.AssertEqual("empty-dish.png", post.Image)
	tu.AssertEqual(1200, post.LikeCount)
	tu.AssertEqual(80, post.RetweetCount)
	tu.AssertEqual(12, post.BookmarkCount)
	tu.AssertEqual(10000, post.Impressions)
	tu.AssertEqual(
		"2024-04-12 11:37:46",
		post.CreatedAt.Format(constants.TIME_LAYOUT))
	tu.AssertEqual(
		"2024-05-18 08:02:13",
		post.UpdatedAt.Format(constants.TIME_LAYOUT))
}

func TestPostNew(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)

		user := NewUserController(db)
		post := NewPostController(db)

		userInput := dtypes.UserInput{
			Username:    "esteban",
			Email:       "estecat42069@yahoo.com",
			Password:    "password",
			DisplayName: "Bubba",
		}

		user.Set(nil, userInput)
		user.Create("password")

		postInput := dtypes.PostInput{
			UserID:  user.ID(),
			Content: "Cats are cool",
			Image:   "teef.jpg",
		}

		beforeAction := time.Now().UTC().Add(-1 * time.Minute)
		post.New(postInput)
		afterAction := time.Now().UTC().Add(time.Minute)

		tu.AssertEqual(1, post.ID)
		tu.AssertEqual(user.ID(), post.UserID)
		tu.AssertEqual("Cats are cool", post.Content)
		tu.AssertEqual("teef.jpg", post.Image)
		tu.AssertEqual("esteban", post.Author.Username)
		tu.AssertEqual("Bubba", post.Author.DisplayName)
		tu.AssertEqual("", post.Author.Avatar)
		tu.AssertEqual(0, post.LikeCount)
		tu.AssertEqual(0, post.RetweetCount)
		tu.AssertEqual(0, post.BookmarkCount)
		tu.AssertEqual(0, post.Impressions)
		tu.AssertTrue(post.CreatedAt.After(beforeAction))
		tu.AssertTrue(post.CreatedAt.Before(afterAction))
		tu.AssertTrue(post.UpdatedAt.After(beforeAction))
		tu.AssertTrue(post.UpdatedAt.Before(afterAction))
	})
}

func TestPostByID(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		before := timestamp.Add(-1 * time.Minute)
		after := timestamp.Add(time.Minute)
		post := NewPostController(db)

		// unknown ID
		err := post.ByID(-1)
		var postNotFoundErr model.PostNotFoundError
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(errors.As(err, &postNotFoundErr))

		err = post.ByID(1)
		tu.AssertErrorNil(err)
		tu.AssertEqual(3, post.UserID)
		tu.AssertEqual(
			"Diane! I'm holding in my hand a small box of chocolate bunnies.",
			post.Content,
		)
		tu.AssertEqual("chocolate-bunnies.png", post.Image)
		tu.AssertEqual(0, post.LikeCount)
		tu.AssertEqual(0, post.RetweetCount)
		tu.AssertEqual(0, post.BookmarkCount)
		tu.AssertEqual(0, post.Impressions)
		tu.AssertTrue(post.CreatedAt.After(before))
		tu.AssertTrue(post.CreatedAt.Before(after))
		tu.AssertTrue(post.UpdatedAt.After(before))
		tu.AssertTrue(post.UpdatedAt.Before(after))
		tu.AssertEqual("dalecooper", post.Author.Username)
		tu.AssertEqual("Coffee Fre@k", post.Author.DisplayName)
		tu.AssertEqual("", post.Author.Avatar)
	})
}

func TestPostNewUserDoesNotExist(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		post := NewPostController(db)
		postInput := dtypes.PostInput{
			UserID:  42069,
			Content: "Some content",
			Image:   "dags.jpg",
		}

		err := post.New(postInput)
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(dbutils.IsConstraintError(err))
	})
}

func TestPostLike(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		user1 := NewUserController(db)
		user2 := NewUserController(db)
		user3 := NewUserController(db)
		user4 := NewUserController(db)
		post := NewPostController(db)
		user1.ByID(1)
		user2.ByID(2)
		user3.ByID(3)
		user4.ByID(4)

		err := post.Like(user1.ID())
		tu.AssertErrorNotNil(err)
		tu.AssertEqual(
			"Post.Like(): missing required postID in post controller",
			err.Error(),
		)

		post.ByID(1)

		err = post.Like(user1.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(1, post.LikeCount)

		err = post.Like(user2.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(2, post.LikeCount)

		err = post.Like(user3.ID())
		tu.AssertErrorNil(err)
		err = post.Like(user4.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(4, post.LikeCount)

		// user4 likes again, no error expected and count stays the same
		err = post.Like(user4.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(4, post.LikeCount)
	})
}

func TestPostUnlike(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		user1 := NewUserController(db)
		user2 := NewUserController(db)
		user3 := NewUserController(db)
		user4 := NewUserController(db)
		post := NewPostController(db)
		user1.ByID(1)
		user2.ByID(2)
		user3.ByID(3)
		user4.ByID(4)

		err := post.Unlike(user1.ID())
		tu.AssertErrorNotNil(err)
		tu.AssertEqual(
			"Post.Unlike(): missing required postID in post controller",
			err.Error(),
		)

		post.ByID(1)
		post.Like(user1.ID())
		post.Like(user2.ID())
		post.Like(user3.ID())
		post.Like(user4.ID())

		err = post.Unlike(user4.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(3, post.LikeCount)

		err = post.Unlike(user4.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(3, post.LikeCount)

		err = post.Unlike(user3.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(2, post.LikeCount)

		err = post.Unlike(user2.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(1, post.LikeCount)

		err = post.Unlike(user1.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(0, post.LikeCount)

		err = post.Unlike(user1.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(0, post.LikeCount)
	})
}

func TestPostRetweet(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		user1 := NewUserController(db)
		user2 := NewUserController(db)
		user3 := NewUserController(db)
		user4 := NewUserController(db)
		post := NewPostController(db)
		user1.ByID(1)
		user2.ByID(2)
		user3.ByID(3)
		user4.ByID(4)

		err := post.Retweet(user1.ID())
		tu.AssertErrorNotNil(err)
		tu.AssertEqual(
			"Post.Retweet(): missing required postID in post controller",
			err.Error(),
		)

		post.ByID(1)

		err = post.Retweet(user1.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(1, post.RetweetCount)

		err = post.Retweet(user2.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(2, post.RetweetCount)

		err = post.Retweet(user3.ID())
		tu.AssertErrorNil(err)

		err = post.Retweet(user4.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(4, post.RetweetCount)

		// user4 retweets again, no error expected and count stays the same
		err = post.Retweet(user4.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(4, post.RetweetCount)
	})
}

func TestPostUnRetweet(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		user1 := NewUserController(db)
		user2 := NewUserController(db)
		user3 := NewUserController(db)
		user4 := NewUserController(db)
		post := NewPostController(db)
		user1.ByID(1)
		user2.ByID(2)
		user3.ByID(3)
		user4.ByID(4)

		err := post.UnRetweet(user1.ID())
		tu.AssertErrorNotNil(err)
		tu.AssertEqual(
			"Post.UnRetweet(): missing required postID in post controller",
			err.Error(),
		)

		post.ByID(1)
		post.Retweet(user1.ID())
		post.Retweet(user2.ID())
		post.Retweet(user3.ID())
		post.Retweet(user4.ID())

		err = post.UnRetweet(user4.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(3, post.RetweetCount)

		err = post.UnRetweet(user4.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(3, post.RetweetCount)

		err = post.UnRetweet(user3.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(2, post.RetweetCount)

		err = post.UnRetweet(user2.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(1, post.RetweetCount)

		err = post.UnRetweet(user1.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(0, post.RetweetCount)

		err = post.UnRetweet(user1.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(0, post.RetweetCount)
	})
}

func TestPostBookmark(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		user1 := NewUserController(db)
		user2 := NewUserController(db)
		user3 := NewUserController(db)
		user4 := NewUserController(db)
		post := NewPostController(db)
		user1.ByID(1)
		user2.ByID(2)
		user3.ByID(3)
		user4.ByID(4)

		err := post.Bookmark(user1.ID())
		tu.AssertErrorNotNil(err)
		tu.AssertEqual(
			"Post.Bookmark(): missing required postID in post controller",
			err.Error(),
		)

		post.ByID(1)

		err = post.Bookmark(user1.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(1, post.BookmarkCount)

		err = post.Bookmark(user2.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(2, post.BookmarkCount)

		err = post.Bookmark(user3.ID())
		tu.AssertErrorNil(err)

		err = post.Bookmark(user4.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(4, post.BookmarkCount)

		// user4 retweets again, no error expected and count stays the same
		err = post.Bookmark(user4.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(4, post.BookmarkCount)
	})
}

func TestPostUnBookmark(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		user1 := NewUserController(db)
		user2 := NewUserController(db)
		user3 := NewUserController(db)
		user4 := NewUserController(db)
		post := NewPostController(db)
		user1.ByID(1)
		user2.ByID(2)
		user3.ByID(3)
		user4.ByID(4)

		err := post.UnBookmark(user1.ID())
		tu.AssertErrorNotNil(err)
		tu.AssertEqual(
			"Post.UnBookmark(): missing required postID in post controller",
			err.Error(),
		)

		post.ByID(1)
		post.Bookmark(user1.ID())
		post.Bookmark(user2.ID())
		post.Bookmark(user3.ID())
		post.Bookmark(user4.ID())

		err = post.UnBookmark(user4.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(3, post.BookmarkCount)

		err = post.UnBookmark(user4.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(3, post.BookmarkCount)

		err = post.UnBookmark(user3.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(2, post.BookmarkCount)

		err = post.UnBookmark(user2.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(1, post.BookmarkCount)

		err = post.UnBookmark(user1.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(0, post.BookmarkCount)

		err = post.UnBookmark(user1.ID())
		tu.AssertErrorNil(err)
		tu.AssertEqual(0, post.BookmarkCount)
	})
}

func TestPostSync(t *testing.T) {
	testutil.WithTestData(t, func(db *sql.DB, timestamp time.Time) {
		tu := testutil.NewTestUtil(t)
		post := NewPostController(db)

		err := post.Sync()
		tu.AssertErrorNotNil(err)
		tu.AssertEqual(
			"Post.Sync(): required postID not set in post controller",
			err.Error(),
		)

		post.ID = -1
		err = post.Sync()
		var postNotFoundErr model.PostNotFoundError
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(errors.As(err, &postNotFoundErr))

		post.ByID(4)
		db.Exec("UPDATE Post SET content = 'wazzup!!!' WHERE id = 4;")
		tu.AssertEqual(
			"Nothing beats a damn fine cup of coffee in the morning.",
			post.Content,
		)

		err = post.Sync()
		tu.AssertErrorNil(err)
		tu.AssertEqual("wazzup!!!", post.Content)
	})
}
