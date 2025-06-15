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
