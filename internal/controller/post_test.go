package controller

import (
	"database/sql"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/dbutils"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/testutil"
)

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
	testutil.WithTestData(t, func(db *sql.DB) {
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
	testutil.WithTestData(t, func(db *sql.DB) {
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
