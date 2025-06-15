package model

import (
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/marcusprice/twitter-clone/internal/dbutils"
	"github.com/marcusprice/twitter-clone/internal/dtypes"
	"github.com/marcusprice/twitter-clone/internal/testutil"
	"github.com/marcusprice/twitter-clone/internal/util"
)

func TestPostNew(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		postModel := NewPostModel(db)
		userInput := dtypes.UserInput{
			Email:       "estecat42069@yahoo.com",
			Username:    "estecat",
			FirstName:   "esteban",
			LastName:    "price",
			DisplayName: "Bubba",
			Password:    "password",
		}
		userID := insertUser(userInput, db)

		postInput := dtypes.PostInput{
			UserID:  userID,
			Content: "Cats are cool",
			Image:   "cool-cats.png",
		}

		beforeAction := time.Now().UTC().Add(-1 * time.Minute)
		postID, err := postModel.New(postInput)
		afterAction := time.Now().UTC().Add(time.Minute)

		postData := queryPost(postID, db)
		createdAt := util.ParseTime(postData.CreatedAt)
		updatedAt := util.ParseTime(postData.UpdatedAt)

		tu.AssertErrorNil(err)
		tu.AssertEqual("Cats are cool", postData.Content)
		tu.AssertEqual("cool-cats.png", postData.Image)
		tu.AssertEqual(1, postID)
		tu.AssertEqual(postID, postData.ID)
		tu.AssertEqual(userID, postData.UserID)
		tu.AssertEqual(0, postData.LikeCount)
		tu.AssertEqual(0, postData.RetweetCount)
		tu.AssertEqual(0, postData.BookmarkCount)
		tu.AssertEqual(0, postData.Impressions)
		tu.AssertTrue(createdAt.After(beforeAction))
		tu.AssertTrue(createdAt.Before(afterAction))
		tu.AssertTrue(updatedAt.After(beforeAction))
		tu.AssertTrue(updatedAt.Before(afterAction))
	})
}

func TestPostNewUserDoesNotExist(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		postModel := NewPostModel(db)
		postInput := dtypes.PostInput{
			UserID:  42069,
			Content: "Some content",
			Image:   "image.png",
		}
		postID, err := postModel.New(postInput)
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(dbutils.IsConstraintError(err))
		tu.AssertTrue(strings.Contains(err.Error(), "FOREIGN KEY constraint failed"))
		tu.AssertEqual(-1, postID)
	})
}

func TestPostGetByID(t *testing.T) {
	testutil.WithTestDB(t, func(db *sql.DB) {
		tu := testutil.NewTestUtil(t)
		postModel := PostModel{db}
		userInput := dtypes.UserInput{
			Email:       "estecat42069@yahoo.com",
			Username:    "estecat",
			FirstName:   "esteban",
			LastName:    "price",
			DisplayName: "Bubba",
			Password:    "password",
		}
		beforeAction := time.Now().UTC().Add(-1 * time.Minute)
		userID := insertUser(userInput, db)
		afterAction := time.Now().UTC().Add(time.Minute)

		postInput := dtypes.PostInput{
			UserID:  userID,
			Content: "Cats are cool",
			Image:   "cool-cats.png",
		}
		postID := insertPost(postInput, db)
		postData, err := postModel.GetByID(postID)
		createdAt := util.ParseTime(postData.CreatedAt)
		updatedAt := util.ParseTime(postData.UpdatedAt)

		tu.AssertErrorNil(err)
		tu.AssertEqual("estecat", postData.Author.Username)
		tu.AssertEqual("Bubba", postData.Author.DisplayName)
		tu.AssertEqual("", postData.Author.Avatar)
		tu.AssertEqual(postID, postData.ID)
		tu.AssertEqual("Cats are cool", postData.Content)
		tu.AssertEqual("cool-cats.png", postData.Image)
		tu.AssertEqual(postID, postData.ID)
		tu.AssertEqual(userID, postData.UserID)
		tu.AssertEqual(0, postData.LikeCount)
		tu.AssertEqual(0, postData.RetweetCount)
		tu.AssertEqual(0, postData.BookmarkCount)
		tu.AssertEqual(0, postData.Impressions)
		tu.AssertTrue(createdAt.After(beforeAction))
		tu.AssertTrue(createdAt.Before(afterAction))
		tu.AssertTrue(updatedAt.After(beforeAction))
		tu.AssertTrue(updatedAt.Before(afterAction))

		_, err = postModel.GetByID(42069)
		var postNotFoundError PostNotFoundError
		tu.AssertErrorNotNil(err)
		tu.AssertTrue(errors.As(err, &postNotFoundError))
	})
}

func queryPost(id int, db *sql.DB) PostData {
	query := `
		SELECT
			id,
			user_id,
			content,
			like_count,
			retweet_count,
			bookmark_count,
			impressions,
			image,
			created_at,
			updated_at
		FROM Post
		WHERE id = $1;
	`
	var postID int
	var userID int
	var content string
	var likeCount int
	var retweetCount int
	var bookmarkCount int
	var impressions int
	var image string
	var createdAt string
	var updatedAt string

	db.
		QueryRow(query, id).
		Scan(
			&postID, &userID, &content, &likeCount, &retweetCount,
			&bookmarkCount, &impressions, &image, &createdAt, &updatedAt)

	postData := PostData{
		ID:            postID,
		UserID:        userID,
		Content:       content,
		LikeCount:     likeCount,
		RetweetCount:  retweetCount,
		BookmarkCount: bookmarkCount,
		Impressions:   impressions,
		Image:         image,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	return postData
}

func insertPost(postInput dtypes.PostInput, db *sql.DB) int {
	query := `
		INSERT INTO Post (user_id, content, image)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	var postID int
	db.
		QueryRow(query, postInput.UserID, postInput.Content, postInput.Image).
		Scan(&postID)

	return postID
}
